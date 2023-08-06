package file

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	proto3 "src/proto/controller_client"
	"strconv"
	"strings"
)

type FileHandler struct {
	fileName   string
	fileSize   int64
	dataStream []byte
	checksum   [16]byte
	location   []string
	dir        string

	fragments   []*Fragment
	fragmentMap map[string]*Fragment
}

func (f *FileHandler) SetDir(dir string) {
	f.dir = dir
}

func (f *FileHandler) Dir() string {
	return f.dir
}

func (f *FileHandler) Location() []string {
	return f.location
}

func (f *FileHandler) SetLocation(location []string) {
	f.location = location
}

func (f *FileHandler) SetChecksum(checksum [16]byte) {
	f.checksum = checksum
}

func (f *FileHandler) Fragments() []*Fragment {
	return f.fragments
}

func (f *FileHandler) SetFragments(fragments []*Fragment) {
	f.fragments = fragments
}

func (f *FileHandler) FragmentMap() map[string]*Fragment {
	return f.fragmentMap
}

func (f *FileHandler) SetFragmentMap(fragmentMap map[string]*Fragment) {
	f.fragmentMap = fragmentMap
}

type Fragment struct {
	fragName     string
	fragPosition int
	totalFrags   int
	fragSize     int64
	fragmentData []byte
	fragChecksum [16]byte
	location     []string
}

func (f *Fragment) Location() []string {
	return f.location
}

func (f *Fragment) SetLocation(location []string) {
	f.location = location
}

func (f *Fragment) FragName() string {
	return f.fragName
}

func (f *Fragment) SetFragName(fragName string) {
	f.fragName = fragName
}

func (f *Fragment) FragSize() int64 {
	return f.fragSize
}

func (f *Fragment) SetFragSize(fragSize int64) {
	f.fragSize = fragSize
}

func (f *Fragment) FragmentData() []byte {
	return f.fragmentData
}

func (f *Fragment) SetFragmentData(fragmentData []byte) {
	f.fragmentData = fragmentData
}

func (f *Fragment) FragChecksum() [16]byte {
	return f.fragChecksum
}

func (f *Fragment) SetFragChecksum(fragChecksum [16]byte) {
	f.fragChecksum = fragChecksum
}

/*
func (f *FileHandler) FillFragmentData(fragId string) []byte {

	fragment := f.fragmentMap[fragId]
	fragmentPos := fragment.fragPosition
	fragmentSize := fragment.fragSize

	file, err := os.Open(DIR + f.fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return nil
	}

	var startPos int64
	var endPos int64

	if fragment.fragPosition == fragment.totalFrags-1 {
		startPos = fileInfo.Size() - fragmentSize
		endPos = fileInfo.Size()

	} else {
		startPos = int64(fragmentPos) * fragmentSize
		endPos = int64(fragmentPos+1) * fragmentSize
		if endPos > fileInfo.Size() {
			endPos = fileInfo.Size()
		}
	}

	buffer := make([]byte, endPos-startPos)
	_, err = file.ReadAt(buffer, startPos)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	fragment.fragChecksum = md5.Sum(buffer)
	return buffer

}
*/
func (f *FileHandler) FillFragmentData(fragId string) []byte {

	fragment := f.fragmentMap[fragId]
	fragmentPos := fragment.fragPosition
	fragmentSize := fragment.fragSize

	file, err := os.Open(f.dir + f.fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var startPos int64
	var endPos int64
	var startPosFile *os.File
	if fragmentPos == 0 {
		startPos = 0
		//create a file to store the end position
		startPosFile, err = os.Create("startPos.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return nil
		}
		defer startPosFile.Close()

	} else {
		startPosFile, err = os.Open("startPos.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil

		}
		defer startPosFile.Close()

		scanner := bufio.NewScanner(startPosFile)
		scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)
		scanner.Scan()
		startPos, err = strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil

		}
	}

	_, err = file.Seek(startPos, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return nil
	}
	var buf bytes.Buffer
	scanner := bufio.NewScanner(file)

	scanner.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if int(fragmentSize) < buf.Len()+len(line) {
			if int(fragmentSize) > buf.Len() {
				buf.Write(line)
				if err := scanner.Err(); err != nil {
					fmt.Println("Error scanning file:", err)
					return nil
				}
				buf.WriteByte('\n')
				endPos = startPos + int64(buf.Len())

			} else {
				break
			}

		} else {
			buf.Write(line)
			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanning file:", err)
				return nil
			}
			buf.WriteByte('\n')
			endPos = startPos + int64(buf.Len())

		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return nil
	}

	startPosFile, err = os.Create("startPos.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	defer startPosFile.Close()

	_, err = startPosFile.WriteString(strconv.FormatInt(endPos, 10))
	if err != nil {
		fmt.Println("Error writing file:", err)
		return nil
	}

	buffer := buf.Bytes()
	fragment.fragChecksum = md5.Sum(buffer)

	return buffer

}

func (f *FileHandler) GetFileName(fragID string) (int, error) {

	lastUnderscore := strings.LastIndex(fragID, "_")
	if lastUnderscore == -1 {
		// If there is no underscore, return the whole string
		return 0, fmt.Errorf("no underscore in fragment name")
	}
	// Otherwise, return the substring after the last underscore, excluding the underscore
	return strconv.Atoi(fragID[lastUnderscore+1:])
}

func (f *FileHandler) SetFragmentLayout(fragments []proto3.FragmentInfo) {

	f.fragmentMap = make(map[string]*Fragment)
	for _, fragment := range fragments {
		frag := &Fragment{
			fragName: fragment.FragmentId,
			fragSize: fragment.Size,
			location: make([]string, 0),
		}
		frag.fragPosition, _ = f.GetFileName(fragment.FragmentId)
		frag.totalFrags = len(fragments)
		for _, location := range fragment.StorageNodes {
			frag.location = append(frag.location, location.Host)
		}
		f.fragments = append(f.fragments, frag)
	}

	// Sort the fragments based on their fragment number
	sort.Slice(f.fragments, func(i, j int) bool {
		f1 := f.fragments[i]
		f2 := f.fragments[j]

		// Extract the fragment numbers from the fragment names
		f1Num, _ := strconv.Atoi(strings.Split(f1.fragName, "_")[1])
		f2Num, _ := strconv.Atoi(strings.Split(f2.fragName, "_")[1])

		return f1Num < f2Num
	})

	//TODO: change this to dynamically read file instead of storing in memory
	/*
		dataLen := 0
		for _, fragment := range f.fragments {
			fr := fragment.fragName
			delimiter := "_"
			split := strings.LastIndex(fr, delimiter)
			parts := []string{fr[:split], fr[split+1:]}
			fmt.Println(parts[1]) // Output: index

			fragment.fragmentData = f.dataStream[dataLen : dataLen+int(fragment.fragSize)]
			dataLen += int(fragment.fragSize)
			fragment.fragChecksum = md5.Sum(fragment.fragmentData)

			f.fragmentMap[fragment.fragName] = fragment

		}

	*/
	for _, fragment := range f.fragments {

		f.fragmentMap[fragment.fragName] = fragment

	}

}

func NewFileHandler(fileName string) *FileHandler {
	newFileHandler := &FileHandler{
		fileName: fileName,
	}
	return newFileHandler

}

func (f *FileHandler) SetDataStream(dataStream []byte) {
	f.dataStream = dataStream
}

func (f *FileHandler) SetFileSize(fileSize int64) {
	f.fileSize = fileSize
}

func (f *FileHandler) SetFileName(fileName string) {
	f.fileName = fileName
}

//func (f *FileHandler) SetDirInfo(dirInfo string) {
//	f.dirInfo = dirInfo
//}

func (f *FileHandler) SetCheckSum(checksum [16]byte) {
	f.checksum = checksum
}

func (f *FileHandler) FindAndSetCheckSum() {
	//set md5 checksum
	f.checksum = md5.Sum(f.dataStream)
}

func (f *FileHandler) Checksum() (checksum []byte) {
	return f.checksum[:]
}

func (f *FileHandler) FileName() string {
	return f.fileName

}

func (f *FileHandler) FileSize() int64 {
	return f.fileSize
}

func (f *FileHandler) DataStream() []byte {
	return f.dataStream
}

//func (f *FileHandler) DirInfo() string {
//	return f.dirInfo
//}

func (f *FileHandler) AppendChunks(chunk []byte) {
	f.dataStream = append(f.dataStream, chunk...)
}

func (f *FileHandler) CompareChecksum(reqChecksum []byte) (equivalent bool) {

	if bytes.Compare(f.checksum[:], reqChecksum) == 0 {
		equivalent = true
		return
	}
	return

}

func (f *FileHandler) DirCheck() (exists bool) {

	return
}

func (f *FileHandler) CalcFileSize() {

	//check the file size using the fileName and set it
	fileOpen, err := os.Open(f.dir + f.fileName)
	if err != nil {
		return
	}

	defer fileOpen.Close()

	fileInfo, _ := fileOpen.Stat()
	f.fileSize = int64(fileInfo.Size())

}

func (f *FileHandler) FileCheck() (exists bool, err error) {
	fmt.Println("file: ", f.fileName)

	if _, e := os.Stat(f.dir + f.fileName); e == nil {
		exists = true
	} else {
		//fmt.Println(err.Error())
		err = e
	}
	return
}

func (f *FileHandler) StorageCheck() (enough bool, err error) {
	//var stat syscall.Statfs_t
	//e := syscall.Statfs(f.fileName, &stat)
	//
	//if e != nil {
	//	err = e
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//available := stat.Bavail * uint64(stat.Bsize)
	//
	//if available < uint64(f.fileSize) {
	//	return
	//}

	enough = true

	return
}

func (f *FileHandler) ExtractFileName() (fileName string) {
	//Take the last part of the path

	fileName = f.fileName[strings.LastIndex(f.fileName, "/")+1:]
	//f.fileName = fileName

	return
}

//func (f *FileHandler) ReadFile() (file2 []byte, err error) {
//
//	//fileName := DIR + f.fileName
//	fileName := DIR + f.fileName
//	//fmt.Println("FILE NAME: ", fileName)
//
//	// define chunk size
//	const chunkSize = 1024 * 1024 // 1MB
//
//	// open the file
//	file, err := os.Open(fileName)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	// get file size
//	fileInfo, _ := file.Stat()
//	fileSize := fileInfo.Size()
//
//	// calculate number of chunks
//	numChunks := int(math.Ceil(float64(fileSize) / float64(chunkSize)))
//
//	// create a channel to communicate between go routines
//	chunkCh := make(chan struct {
//		index int
//		data  []byte
//	}, numChunks)
//
//	// create a wait group to wait for all go routines to complete
//	wg := sync.WaitGroup{}
//	wg.Add(numChunks)
//
//	// start go routines to read chunks of the file concurrently
//	for i := 0; i < numChunks; i++ {
//		go func(chunkIndex int) {
//			defer wg.Done()
//
//			// calculate offset and length of the chunk
//			offset := int64(chunkIndex * chunkSize)
//			length := int(math.Min(float64(chunkSize), float64(fileSize-offset)))
//
//			// read the chunk
//			chunk := make([]byte, length)
//			_, err := file.ReadAt(chunk, offset)
//			if err != nil && err != io.EOF {
//				log.Printf("Error reading chunk %d: %s\n", chunkIndex, err.Error())
//				return
//			}
//
//			// send the chunk through the channel
//			chunkCh <- struct {
//				index int
//				data  []byte
//			}{index: chunkIndex, data: chunk}
//		}(i)
//	}
//
//	// create a buffer and wait for chunks to be received from the channel
//	buffer := make([][]byte, numChunks)
//	go func() {
//		for chunk := range chunkCh {
//			buffer[chunk.index] = chunk.data
//		}
//	}()
//
//	// wait for all go routines to finish
//	wg.Wait()
//	close(chunkCh)
//
//	// combine chunks into a single byte slice
//	result := make([]byte, 0, fileSize)
//	for _, chunk := range buffer {
//		result = append(result, chunk...)
//	}
//	f.SetDataStream(result)
//
//	return result, nil
//}

/*
func (f *FileHandler) ReadFile() (file []byte, err error) {
	fileName := DIR + f.fileName
	data, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Couldn't open file", fileName)
	}

	defer data.Close()

	scanner := bufio.NewScanner(data)

	const maxCapacity = 11000000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	file = make([]byte, 0)
	for scanner.Scan() {
		file = append(file, scanner.Bytes()...)
	}
	f.SetDataStream(file)

	return

}
*/

func (f *FileHandler) ReadFile() (file []byte, err error) {

	fileName := f.dir + f.fileName
	fmt.Println("FILE NAME: ", fileName)
	fileOpen, err := os.Open(fileName)
	if err != nil {

		return
	}

	defer fileOpen.Close()

	fileInfo, _ := fileOpen.Stat()
	fmt.Println("FILE SIZE: ", fileInfo.Size())
	f.fileSize = fileInfo.Size()

	file, err = ioutil.ReadAll(fileOpen)
	if err != nil {
		return
	}
	f.SetDataStream(file)

	return

}

func (f *FileHandler) WriteFile() (err error) {

	var dir string

	//dir = DIR + f.fileName
	dir = f.dir + f.fileName

	file, err := os.Create(dir)
	if err != nil {
		fmt.Println("Error Writing the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(f.dataStream)
	if err != nil {
		fmt.Println("Error Writing the file")
		fmt.Println(err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error Writing the file")
		fmt.Println(err)
		return
	}

	fmt.Println("File written successfully")

	return
}

func (f *FileHandler) ChecksumOnDisk() (err error) {

	//write the checksum to a file on disk. File name will be the checksum

	//dir = dir + "checksum"
	path := f.dir + f.fileName + ".checksum"

	//TODO: check if checksum alredy exists
	output := fmt.Sprintf("%s %x\n", f.fileName, f.Checksum())

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error Writing the file")
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(output); err != nil {
		fmt.Println("Error Writing the file")
		panic(err)
	}
	return

}

func (f *FileHandler) ReadCheckSum() (files []string, err error) {

	var dir string

	dir = f.dir + "checksum"

	file, err := os.Open(dir)
	if err != nil {
		fmt.Println("Error Reading the file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//split the string and get the file name:
		fileName := strings.Split(scanner.Text(), " ")
		//fmt.Println(fileName[0])
		files = append(files, fileName[0])
	}

	return

	//checksum stored in filename.checksum

	//for _, file := range files {
	//
	//	dir := DIR + file + ".checksum"
	//	openFile, err := os.Open(dir)
	//	if err != nil {
	//		fmt.Println("Error Reading the file")
	//	}
	//	defer openFile.Close()
	//
	//	scanner := bufio.NewScanner(openFile)
	//	for scanner.Scan() {
	//		//check stored as checksum
	//		if strings.Contains(scanner.Text(), "checksum") {
	//
	//		}
	//	}
	//
	//}
	//
	//return
}

func (f *FileHandler) ValidateChecksumFromFile(fileName string) (valid bool, err error) {

	var dir string

	dir = f.dir + fileName + ".checksum"

	file, err := os.Open(dir)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), hex.EncodeToString(f.Checksum()[:])) {
			valid = true
			return
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return

}
