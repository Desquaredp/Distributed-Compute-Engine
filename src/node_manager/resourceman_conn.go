package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"plugin"
	"sort"
	messagesNodeManager "src/messages/nodeman_nodeman"
	resourceManagerMessages "src/messages/resourceman_nodeman"
	proto3NodeManager "src/proto/nodeman_nodeman"
	resourceManagerProto3 "src/proto/resourceman_nodeman"
	"src/types"
	"strconv"
	"strings"
)

//var nm.dir = "../../../../../../bigdata/students/dmistry4/"

func (nm *NodeManager) AcceptResourceManagerConnections(listener net.Listener) {
	for {

		conn, err := listener.Accept()
		if err != nil {
			nm.logger.Fatal("Error accepting connection")
		} else {

			nm.logger.Info("Resource Manager connected")
			msgHandler := resourceManagerMessages.NewMessageHandler(conn)
			go nm.HandleResourceManagerConnection(msgHandler)
		}

	}
}

func (nm *NodeManager) HandleResourceManagerConnection(msgHandler *resourceManagerMessages.MessageHandler) {

	//defer msgHandler.Close()
	proto := resourceManagerProto3.NewProtoHandler(msgHandler, nm.logger)
	defer proto.MsgHandler().Close()

	for {
		wrapper, _ := proto.MsgHandler().ResourceManagerRequestReceive()

		switch wrapper.Task.(type) {
		default:
			req := proto.HandleResourceManagerRequest(wrapper)
			nm.logger.Info("Received task request from the resource manager.")

			switch req.(type) {
			case *resourceManagerProto3.SetupRequest:
				request := req.(*resourceManagerProto3.SetupRequest)
				nm.logger.Info("Processing setup request")
				nm.processSetupRequest(request, proto)

			case *resourceManagerProto3.FragmentRequest:
				request := req.(*resourceManagerProto3.FragmentRequest)
				nm.logger.Info("Processing fragment request")
				nm.processFragmentRequest(request, proto)

			case *resourceManagerProto3.KeyDistributionRequest:
				request := req.(*resourceManagerProto3.KeyDistributionRequest)
				nm.logger.Info("Processing key distribution request")
				nm.processKeyDistributionRequest(request, proto)

			case *resourceManagerProto3.ReduceTaskRequest:
				request := req.(*resourceManagerProto3.ReduceTaskRequest)
				nm.logger.Info("Processing reduce request")
				nm.processReduceComplete(request, proto)

			}

		case nil:
			return

		}

	}

}

func (nm *NodeManager) processReduceComplete(request *resourceManagerProto3.ReduceTaskRequest, proto *resourceManagerProto3.ProtoHandler) {
	nm.logger.Info("Received reduce request from the resource manager.")
	nm.logger.Info("Number of reduce operations: " + strconv.Itoa(len(request.TaskID)))

	reducerEncoderMap := make(map[int32]*ReducerEncoder)

	for _, taskID := range request.TaskID {

		reducerEncoderMap[taskID] = &ReducerEncoder{
			reader:   nil,
			encoder:  nil,
			OpenFile: nil,
		}
	}

	mr := nm.loadReducerPlugin(request.InputFileName)
	for _, taskID := range request.TaskID {
		nm.logger.Info("TaskID: " + strconv.Itoa(int(taskID)))
		reducerFiles := nm.findFilesToReduce(taskID, request.InputFileName)
		nm.logger.Info("Number of files to reduce: " + strconv.Itoa(len(reducerFiles)))
		nm.logger.Info("Files to reduce: " + strings.Join(reducerFiles, ", "))
		queue, err := nm.populateQueue(reducerFiles)
		if err != nil {
			nm.logger.Error("Error populating queue")

			//TODO: send error to resource manager
			proto.HandleReduceTaskUpdate(taskID, false)
			return
		}
		nm.logger.Info("Starting reduce operation")
		for queue.Size() > 0 {

			dequeuedFiles := make([]*File, 0)
			dequeuedFile, boolVal := queue.Dequeue()
			if !boolVal {
				nm.logger.Error("Error dequeuing file")
				proto.HandleReduceTaskUpdate(taskID, false)
				return
			}
			dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))

			for queue.Size() > 0 {
				nextFile, boolVal := queue.Peek()
				if !boolVal {
					nm.logger.Error("Error peeking file")
					err = errors.New("Error peeking file")
					//TODO: send error to resource manager
					proto.HandleReduceTaskUpdate(taskID, false)
					return
				}
				if bytes.Compare(dequeuedFile.(*File).CurrKV.Key, nextFile.(*File).CurrKV.Key) == 0 {
					dequeuedFile, boolVal = queue.Dequeue()
					if !boolVal {
						nm.logger.Error("Error dequeuing file")
						err = errors.New("Error dequeuing file")
						//TODO: send error to resource manager
						proto.HandleReduceTaskUpdate(taskID, false)
						return
					}
					dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))
				} else {
					break
				}
			}

			reducerID := taskID
			/*
				Steps:
					- get the key values from the files(done)
					- perform the reduce operation (done)
					- encode the result to the output file (done)

			*/

			kvPair := nm.extractKVPair(dequeuedFiles)

			reduceOperationOutput := mr.Reduce(*kvPair)

			err = nm.encodeToOutputFile(request.OutputFileName, reduceOperationOutput, reducerID, reducerEncoderMap)
			if err != nil {
				nm.logger.Error("Error encoding to output file: " + err.Error())
				proto.HandleReduceTaskUpdate(taskID, false)
				return
			}

			for _, file := range dequeuedFiles {
				file.CurrKV = types.ReduceInputRecord{}
				err = file.Decoder.Decode(&file.CurrKV)
				if err != nil {
					if err == io.EOF {
						//close the file
						file.OpenFile.Close()
						err = nil
						continue
					} else {
						nm.logger.Error("Error decoding file: " + err.Error())
						proto.HandleReduceTaskUpdate(taskID, false)
						return
					}
				}
				queue.Enqueue(file)
			}

		}

		nm.logger.Info("Deleting reducer files")
		err = nm.deleteReducerFiles(reducerFiles)
		if err != nil {
			nm.logger.Error("Error deleting reducer files: " + err.Error())
			proto.HandleReduceTaskUpdate(taskID, false)
			return
		}
		proto.HandleReduceTaskUpdate(taskID, true)

	}

}

func (nm *NodeManager) deleteReducerFiles(files []string) (err error) {

	for _, file := range files {
		err = os.Remove(nm.dir + file)
		if err != nil {
			nm.logger.Error("Error deleting intermediate file: " + err.Error())
			return
		}
	}

	return

}
func (nm *NodeManager) extractKVPair(files []*File) *types.ReduceInputRecord {

	values := make([][]byte, 0)
	for _, file := range files {
		values = append(values, file.CurrKV.Values...)
	}

	kvPair := &types.ReduceInputRecord{
		Key:    files[0].CurrKV.Key,
		Values: values,
	}

	return kvPair
}

func (nm *NodeManager) loadReducerPlugin(inputFileName string) types.MRPluginInterface {

	filePath := nm.dir + inputFileName + "_MR_job" + ".so"

	mod := filePath

	nm.logger.Info("Loading plugin: " + mod)

	p, err := plugin.Open(mod)
	if err != nil {
		nm.logger.Error("Error opening plugin: " + err.Error())
	}

	mrSym, err := p.Lookup("MRPlugin")
	if err != nil {
		nm.logger.Error("Error looking up symbol: " + err.Error())
	}

	var mr types.MRPluginInterface
	mr, ok := mrSym.(types.MRPluginInterface)
	if !ok {
		nm.logger.Error("Error casting to MRPluginInterface")
	}

	return mr
}

/*
func (nm *NodeManager) encodeToOutputFile(outputFileName string, reduceOperationOutput types.ReduceOutputRecord, reducerID int32, encoderMap map[int32]*ReducerEncoder) (err error) {

	if encoderMap[reducerID].reader == nil {
		reducerFile := nm.dir + outputFileName + "_MR_output_" + strconv.Itoa(int(reducerID))
		outputFile, err := os.OpenFile(reducerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			nm.logger.Error("Error opening file: " + err.Error())
			return err
		}

		encoder := gob.NewEncoder(outputFile)

		encoderMap[reducerID].reader = outputFile
		encoderMap[reducerID].encoder = encoder
		err = encoder.Encode(types.ReduceOutputRecord{Key: reduceOperationOutput.Key, Value: reduceOperationOutput.Value})
		if err != nil {
			nm.logger.Error("Error encoding to file: " + err.Error())
			return err
		}

	} else {

		err = encoderMap[reducerID].encoder.Encode(types.ReduceOutputRecord{Key: reduceOperationOutput.Key, Value: reduceOperationOutput.Value})
		if err != nil {
			nm.logger.Error("Error encoding to file: " + err.Error())
			return err
		}
	}

	return

}

*/

func (nm *NodeManager) encodeToOutputFile(outputFileName string, reduceOperationOutput types.ReduceOutputRecord, reducerID int32, encoderMap map[int32]*ReducerEncoder) (err error) {

	reducerFile := nm.dir + outputFileName + "_MR_output_" + strconv.Itoa(int(reducerID))
	outputFile, err := os.OpenFile(reducerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		nm.logger.Error("Error opening file: " + err.Error())
		return err
	}
	defer outputFile.Close()

	//write to file as string with key \t value

	_, err = outputFile.WriteString(string(reduceOperationOutput.Key) + "\t" + string(reduceOperationOutput.Value) + "\n")
	if err != nil {
		nm.logger.Error("Error writing to file: " + err.Error())
		return err
	}

	return
}

func (nm *NodeManager) findFilesToReduce(reducerID int32, name string) (filesToReduce []string) {

	filesToReduce = make([]string, 0)

	f, err := os.Open(nm.dir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "_MR_reducer_"+strconv.Itoa(int(reducerID))) && strings.Contains(file.Name(), name) {
			filesToReduce = append(filesToReduce, file.Name())
		}
	}

	return
}

type ReducerEncoder struct {
	reader   io.Reader
	encoder  *gob.Encoder
	OpenFile *os.File
}

func (nm *NodeManager) processKeyDistributionRequest(request *resourceManagerProto3.KeyDistributionRequest, proto *resourceManagerProto3.ProtoHandler) {
	//Print the request

	nm.logger.Info("Received key distribution request from the resource manager.")
	storageNodes := request.StorageNode
	reducerKeyDistribution := request.ReducerKeyDistribution

	for _, storageNode := range storageNodes {
		nm.logger.Info("Storage Node: " + storageNode.Host)
		nm.logger.Info("ReducerIDs assigned: ")
		for _, reducerID := range storageNode.ReduceTaskIds {
			nm.logger.Info(strconv.Itoa(int(reducerID)))
		}
	}

	for _, reducerKey := range reducerKeyDistribution {
		nm.logger.Info("ReducerID: " + strconv.Itoa(int(reducerKey.ReducerID)))
		nm.logger.Info("Keys assigned: " + strconv.Itoa(len(reducerKey.Keys)))

	}

	//locateReducerMap -> mapping reducerID to the storage node
	//reduceKeysMap -> mapping reducerID to keys

	locateReducerMap := make(map[int32]*resourceManagerProto3.StorageNode)
	reduceKeysMap := make(map[int32][][]byte)

	for _, storageNode := range storageNodes {
		for _, reducerID := range storageNode.ReduceTaskIds {
			locateReducerMap[reducerID] = storageNode
		}
	}

	for _, reducerKey := range reducerKeyDistribution {
		reduceKeysMap[reducerKey.ReducerID] = reducerKey.Keys
	}

	//map keys -> reducerID
	keysToReducerIDs := make(map[string]int32)

	for reducerID, keys := range reduceKeysMap {
		for _, key := range keys {
			keysToReducerIDs[string(key)] = reducerID
		}
	}

	//Create  a  map: reducerID -> struct {reader, encoder}

	reducerEncoderMap := make(map[int32]*ReducerEncoder)

	for reducerID, _ := range locateReducerMap {

		reducerEncoderMap[reducerID] = &ReducerEncoder{
			reader:   nil,
			encoder:  nil,
			OpenFile: nil,
		}

	}

	/*
		Steps:
			- get all the intermediate files from the storage nodes (done)
			- check all the intermediate files and perform look at one key at a time (done)
		    - check where the key belongs append it to a file named after the reducerID (done)
			- send the files to their respective storage nodes
			- once all the files are sent, send a message to the resource manager

	*/

	files := nm.findIntermediateFiles()
	fileName, err := nm.distributeIntermediateFiles(files, reducerEncoderMap, keysToReducerIDs)
	if err != nil {
		nm.logger.Error("Error distributing intermediate files")
		nm.logger.Error(err.Error())
		return
	}

	responseChan := make(chan proto3NodeManager.ResponseInterface)
	done := make(chan bool)

	nm.logger.Info("Sending reducer files to its respective storage nodes")
	reducersSent := nm.dispatchReducerFiles(locateReducerMap, fileName, responseChan)

	var success bool
	if reducersSent != 0 {

		go func() {
			count := 0
			for response := range responseChan {
				nm.logger.Info("Received response from a node for reducer file")

				res := response.(*proto3NodeManager.ReducerFileResponse)
				if res.Verdict {
					count++
					if count == reducersSent {
						nm.logger.Info("All reducer files have been sent to their respective storage nodes")
						done <- true
					}
				} else {

					nm.logger.Error("Error sending reducer file to storage node")
					done <- false
				}

			}

		}()

		success = <-done
	} else {
		success = true
	}

	nm.logger.Info("Deleting intermediate files")
	err = nm.deleteIntermediateFiles(files)
	if err != nil {
		nm.logger.Error("Error deleting intermediate files")
		nm.logger.Error(err.Error())
	}
	nm.logger.Info("Sending message to resource manager")
	proto.HandleShuffleComplete(success)

}

func (nm *NodeManager) deleteIntermediateFiles(files []string) (err error) {

	for _, file := range files {
		err = os.Remove(nm.dir + file)
		if err != nil {
			nm.logger.Error("Error deleting intermediate file: " + err.Error())
			return
		}
	}

	return
}
func (nm *NodeManager) dispatchReducerFiles(locateReducerMap map[int32]*resourceManagerProto3.StorageNode, fileName string, responseChan chan proto3NodeManager.ResponseInterface) (reducersSent int) {

	reducersSent = 0
	for reducerID, storageNode := range locateReducerMap {

		isCurr, err := nm.checkIfCurrentNode(storageNode)
		if err != nil {
			nm.logger.Error("Error checking if the storage node is the current node")
			nm.logger.Error(err.Error())
			return
		}

		if !isCurr {
			reducersSent++
			nm.logger.Info("Sending reducer file to storage node: " + storageNode.Host)
			nm.sendReducerFileToStorageNode(reducerID, storageNode, fileName, responseChan)
		} else {
			nm.logger.Info("Current node is the storage node for reducerID: " + strconv.Itoa(int(reducerID)))
			continue
		}

	}

	return
}

func (nm *NodeManager) checkIfCurrentNode(storageNode *resourceManagerProto3.StorageNode) (isCurrentNode bool, err error) {
	/*
		ip := net.ParseIP(storageNode.Host)
		if ip == nil {
			nm.logger.Error("Invalid IP address")
			return false, errors.New("Invalid IP address")
		}
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			nm.logger.Error("Failed to get interface addresses")
			return false, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.Contains(ip) {
				nm.logger.Info(ip.String() + "is associated with the current host")
				return true, nil
			}
		}

		return false, nil
	*/

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// If the hostname matches "orion01", we are on that machine
	if strings.Contains(hostname, storageNode.Host) {
		nm.logger.Info("Current node is the storage node for reducerID: " + strconv.Itoa(int(storageNode.ReduceTaskIds[0])))
		return true, nil
	} else {

		return false, nil
	}

}

func (nm *NodeManager) sendReducerFileToStorageNode(reducerID int32, storageNode *resourceManagerProto3.StorageNode, fileName string, responseChan chan proto3NodeManager.ResponseInterface) (err error) {

	nm.logger.Info("Sending reducer file to storage node: " + storageNode.Host)
	nm.logger.Info("Dialing NodeManager for node: " + storageNode.Host)
	msgHandler, err := nm.dialNodeManager(storageNode.Host)
	if err != nil {
		nm.logger.Info("Error dialing the node manager")
		nm.logger.Info(err.Error())
		return
	}

	go nm.handleNodeManagerMessages(msgHandler, responseChan)
	protoInstance := proto3NodeManager.NewProtoHandler(msgHandler, nm.logger)
	nm.HandleReducerFileTransfer(reducerID, protoInstance, fileName)

	return
}

func (nm *NodeManager) HandleReducerFileTransfer(reducerID int32, proto *proto3NodeManager.ProtoHandler, name string) (err error) {

	reducerFile := nm.dir + name + "_MR_reducer_" + strconv.Itoa(int(reducerID))
	nm.logger.Info("Reducer file name: " + reducerFile)
	file, err := os.Open(reducerFile)
	if err != nil {
		nm.logger.Error("Error opening reducer file")
		nm.logger.Error(err.Error())
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 15*1024*1024), 15*1024*1024)

	scanner.Split(bufio.ScanBytes)
	var data []byte
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	checksum := md5.Sum(data)

	proto.SendReducerFileTransfer(reducerFile, data, checksum)

	return
}

func (nm *NodeManager) handleNodeManagerMessages(msgHandler *messagesNodeManager.MessageHandler, responseChan chan proto3NodeManager.ResponseInterface) {

	defer msgHandler.Close()
	proto := proto3NodeManager.NewProtoHandler(msgHandler, nm.logger)
	for {

		wrapper, _ := proto.MsgHandler().ClientRequestReceive()

		switch wrapper.Request.(type) {
		default:
			res := proto.HandleNodeManagerResponse(wrapper)
			nm.logger.Info("Received response from node manager: ")

			switch res.(type) {

			case *proto3NodeManager.ReducerFileResponse:
				response := res.(*proto3NodeManager.ReducerFileResponse)
				if response.Verdict {
					nm.logger.Info("Reducer file transfer was successful")
					responseChan <- response

				} else {
					nm.logger.Info("Reducer file transfer was unsuccessful")
					responseChan <- response
				}

			}

		case nil:
			nm.logger.Info("Received nil message")
			return
		}

	}
}

func (nm *NodeManager) dialNodeManager(ip string) (msgHandler *messagesNodeManager.MessageHandler, err error) {
	conn, err := net.Dial("tcp", ip+":"+"23098")
	if err != nil {
		nm.logger.Error("Error dialing the node manager")
		nm.logger.Error(err.Error())
		return
	}

	msgHandler = messagesNodeManager.NewMessageHandler(conn)
	return
}

func (nm *NodeManager) findIntermediateFiles() (intermediateFiles []string) {
	//get all the files in the directory and check if the names contain "_MR_intermediate"

	intermediateFiles = make([]string, 0)

	f, err := os.Open(nm.dir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "_MR_intermediate") {
			intermediateFiles = append(intermediateFiles, file.Name())
		}
	}

	return
}

type File struct {
	Reader   io.Reader
	Decoder  *gob.Decoder
	CurrKV   types.ReduceInputRecord
	OpenFile *os.File
}

func byBytes(a, b interface{}) int {
	fileA := a.(*File)
	fileB := b.(*File)
	return bytes.Compare(fileA.CurrKV.Key, fileB.CurrKV.Key)

	//return utils.ByteComparator(fileA.CurrKV.Key, fileB.CurrKV.Key)
}

func (nm *NodeManager) distributeIntermediateFiles(intermediateFiles []string, encoderMap map[int32]*ReducerEncoder, keysToReducerIDs map[string]int32) (fileName string, err error) {
	//extrace file names from the intermediate files

	fileNameIndex := strings.Index(intermediateFiles[0], "_MR_intermediate")
	fileName = intermediateFiles[0][:fileNameIndex]
	nm.logger.Info("File name: " + fileName)

	queue, err := nm.populateQueue(intermediateFiles)
	if err != nil {
		nm.logger.Error("Error populating queue: " + err.Error())
		return fileName, err
	}

	for queue.Size() > 0 {

		dequeuedFiles := make([]*File, 0)
		dequeuedFile, boolVal := queue.Dequeue()
		if !boolVal {
			nm.logger.Error("Error dequeuing file")
			return
		}
		dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))

		for queue.Size() > 0 {
			nextFile, boolVal := queue.Peek()
			if !boolVal {
				nm.logger.Error("Error peeking file")
				err = errors.New("Error peeking file")
				return fileName, err
			}
			if bytes.Compare(dequeuedFile.(*File).CurrKV.Key, nextFile.(*File).CurrKV.Key) == 0 {
				dequeuedFile, boolVal = queue.Dequeue()
				if !boolVal {
					nm.logger.Error("Error dequeuing file")
					err = errors.New("Error dequeuing file")
					return
				}
				dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))
			} else {
				break
			}
		}

		reducerID := nm.getReducerID(dequeuedFile.(*File).CurrKV.Key, keysToReducerIDs)
		if reducerID == -1 {

			nm.logger.Error("Error getting reducerID for : " + string(dequeuedFile.(*File).CurrKV.Key))

			err = errors.New("Error getting reducerID")
			return
		}

		err = nm.encodeToReducerFile(fileName, dequeuedFiles, reducerID, encoderMap)
		if err != nil {
			nm.logger.Error("Error encoding to reducer file: " + err.Error())
			return fileName, err
		}

		for _, file := range dequeuedFiles {
			file.CurrKV = types.ReduceInputRecord{}
			err = file.Decoder.Decode(&file.CurrKV)
			if err != nil {
				if err == io.EOF {
					file.OpenFile.Close()
					err = nil
					continue
				} else {
					nm.logger.Error("Error decoding file: " + err.Error())
					return fileName, err
				}
			}
			queue.Enqueue(file)
		}

	}

	//nm.readPopulatedFile(fileName)

	nm.logger.Info("Finished distributing intermediate files")

	return
}

func (nm *NodeManager) readPopulatedFile(fileName string) (err error) {

	reducerFile := nm.dir + fileName + "_MR_reducer_" + strconv.Itoa(int(0))

	outputFile, err := os.OpenFile(reducerFile, os.O_RDONLY, 0644)

	if err != nil {
		nm.logger.Error("Error opening file: " + err.Error())
		return err
	}

	defer outputFile.Close()

	decoder := gob.NewDecoder(outputFile)

	for {

		var reduceInputRecord types.ReduceInputRecord
		err = decoder.Decode(&reduceInputRecord)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				nm.logger.Error("Error decoding file: " + err.Error())
				return err
			}
		}

	}

	return

}

func (nm *NodeManager) encodeToReducerFile(fileName string, dequeuedFiles []*File, reducerID int32, encoderMap map[int32]*ReducerEncoder) (err error) {

	values := make([][]byte, 0)
	for _, file := range dequeuedFiles {
		values = append(values, file.CurrKV.Values...)
	}
	if encoderMap[reducerID].reader == nil {

		reducerFile := nm.dir + fileName + "_MR_reducer_" + strconv.Itoa(int(reducerID))
		outputFile, err := os.OpenFile(reducerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			nm.logger.Error("Error opening file: " + err.Error())
			return err
		}

		encoderMap[reducerID].reader = outputFile
		encoderMap[reducerID].encoder = gob.NewEncoder(outputFile)

		encoder := encoderMap[reducerID].encoder

		err = encoder.Encode(types.ReduceInputRecord{Key: dequeuedFiles[0].CurrKV.Key, Values: values})
		if err != nil {
			nm.logger.Error("Error encoding to file: " + err.Error())
			return err
		}

	} else {

		encoder := encoderMap[reducerID].encoder

		err = encoder.Encode(types.ReduceInputRecord{Key: dequeuedFiles[0].CurrKV.Key, Values: values})
		if err != nil {
			nm.logger.Error("Error encoding to file: " + err.Error())
			return err
		}
	}

	return
}

func (nm *NodeManager) getReducerID(key []byte, keysToReducerIDs map[string]int32) (reducerID int32) {

	keyString := string(key)
	reducerID, ok := keysToReducerIDs[keyString]
	if ok {
		return reducerID
	} else {
		nm.logger.Error("Error getting reducerID")
	}

	return -1
}

func (nm *NodeManager) populateQueue(intermediateFiles []string) (queue *pq.Queue, err error) {
	queue = pq.NewWith(byBytes)

	nm.logger.Info("Populating queue")
	for _, filename := range intermediateFiles {
		inputFile, err := os.Open(nm.dir + filename)
		if err != nil {
			nm.logger.Error("Error opening file: " + err.Error())
			return nil, err
		}

		decoder := gob.NewDecoder(inputFile)
		CurrKV := types.ReduceInputRecord{}
		err = decoder.Decode(&CurrKV)
		if err != nil {
			//Cannot have io.EOF here
			nm.logger.Error("Error decoding file: " + err.Error())
			return nil, err
		}

		file := &File{
			Reader:   inputFile,
			Decoder:  decoder,
			CurrKV:   CurrKV,
			OpenFile: inputFile,
		}
		queue.Enqueue(file)
	}

	return

}

func (nm *NodeManager) processFragmentRequest(response *resourceManagerProto3.FragmentRequest, proto *resourceManagerProto3.ProtoHandler) {

	fragName := response.Fragment
	fileName := response.FileName

	/*
		Steps:
		- Check if file exists
			- If file exists, plug it into the job(done)
			- If file does not exist, send error message to resource manager

	*/

	filePath := nm.dir + fileName + "_MR_job" + ".so"
	fragPath := nm.dir + fragName

	// Open the file
	file, err := os.Open(fragPath)
	if err != nil {
		nm.logger.Error("Error opening file: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 15*1024*1024), 15*1024*1024)

	mod := filePath

	nm.logger.Info("Loading plugin: " + mod)

	p, err := plugin.Open(mod)
	if err != nil {
		nm.logger.Error("Error opening plugin: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	mrSym, err := p.Lookup("MRPlugin")
	if err != nil {
		nm.logger.Error("Error looking up symbol: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	var mr types.MRPluginInterface
	mr, ok := mrSym.(types.MRPluginInterface)
	if !ok {
		nm.logger.Error("Error casting to MRPluginInterface")
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	//TODO: Instead of saving the intermediate in a map(will most likely cause memory issues), sort the intermediate and save it to a file

	var output []types.MapOutputRecord
	for scanner.Scan() {
		line := scanner.Text()
		lineOutput := mr.Map([]byte(line))
		output = append(output, lineOutput...)

	}

	err = nm.sortKeys(output)
	if err != nil {
		nm.logger.Error("Error sorting keys: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	keys, err := nm.writeIntermediateFiles(output, fragPath)
	if err != nil {
		nm.logger.Error("Error writing intermediate files: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}
	output = nil

	nm.logger.Info("Map task completed")
	nm.HandleProcessUpdate(proto, fragName, true, keys)

	//TODO: Send message to resource manager that map task is complete

}

/*
func (nm *NodeManager) processFragmentRequest(response *resourceManagerProto3.FragmentRequest, proto *resourceManagerProto3.ProtoHandler) {

	fragName := response.Fragment
	fileName := response.FileName
	filePath := nm.dir + fileName + "_MR_job" + ".so"
	fragPath := nm.dir + fragName

	// Open the file
	file, err := os.Open(fragPath)
	if err != nil {
		nm.logger.Error("Error opening file: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 15*1024*1024), 15*1024*1024)

	mod := filePath

	nm.logger.Info("Loading plugin: " + mod)

	p, err := plugin.Open(mod)
	if err != nil {
		nm.logger.Error("Error opening plugin: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	mrSym, err := p.Lookup("MRPlugin")
	if err != nil {
		nm.logger.Error("Error looking up symbol: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	var mr types.MRPluginInterface
	mr, ok := mrSym.(types.MRPluginInterface)
	if !ok {
		nm.logger.Error("Error casting to MRPluginInterface")
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	//TODO: Instead of saving the intermediate in a map(will most likely cause memory issues), sort the intermediate and save it to a file

	chunkSize := 10 * 1024 * 1024 //10MB
	index := 0
	chunkFileNames := make([]string, 0)
	var output []types.MapOutputRecord
	for scanner.Scan() {
		line := scanner.Text()
		lineOutput := mr.Map([]byte(line))
		output = append(output, lineOutput...)

		if len(output) >= chunkSize {
			err = nm.sortKeys(output)
			if err != nil {
				nm.logger.Error("Error sorting keys: " + err.Error())
				nm.HandleProcessUpdate(proto, fragName, false, nil)
			}

			fileChunk := fragPath + "_" + strconv.Itoa(index)
			_, err := nm.writeIntermediateFiles(output, fileChunk)
			if err != nil {
				nm.logger.Error("Error writing intermediate files: " + err.Error())
				nm.HandleProcessUpdate(proto, fragName, false, nil)
			}

			chunkFileNames = append(chunkFileNames, fileChunk)
			output = make([]types.MapOutputRecord, 0)
			index++

		} else {
			continue
		}

	}

	if len(output) > 0 {

		err = nm.sortKeys(output)
		if err != nil {
			nm.logger.Error("Error sorting keys: " + err.Error())
			nm.HandleProcessUpdate(proto, fragName, false, nil)
		}

		fileChunk := fragPath + "_" + strconv.Itoa(index)
		_, err := nm.writeIntermediateFiles(output, fileChunk)
		if err != nil {
			nm.logger.Error("Error writing intermediate files: " + err.Error())
			nm.HandleProcessUpdate(proto, fragName, false, nil)
		}
		output = make([]types.MapOutputRecord, 0)
	}

	queue, err := nm.populateQueue(chunkFileNames)
	if err != nil {
		nm.logger.Error("Error populating queue: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
	}

	fragFilePath := fragPath + "_MR_intermediate"
	outputFile, err := os.OpenFile(fragFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		nm.logger.Error("Error opening file: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
		return
	}
	defer outputFile.Close()

	encoder := gob.NewEncoder(outputFile)
	keys := make([][]byte, 0)
	for queue.Size() > 0 {

		dequeuedFiles := make([]*File, 0)
		dequeuedFile, boolVal := queue.Dequeue()
		if !boolVal {
			nm.logger.Error("Error dequeuing file")
			return
		}
		dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))

		for queue.Size() > 0 {
			nextFile, boolVal := queue.Peek()
			if !boolVal {
				nm.logger.Error("Error peeking file")
				err = errors.New("Error peeking file")

				nm.HandleProcessUpdate(proto, fragName, false, nil)
				return
			}
			if bytes.Compare(dequeuedFile.(*File).CurrKV.Key, nextFile.(*File).CurrKV.Key) == 0 {
				dequeuedFile, boolVal = queue.Dequeue()
				if !boolVal {
					nm.logger.Error("Error dequeuing file")
					err = errors.New("Error dequeuing file")
					nm.HandleProcessUpdate(proto, fragName, false, nil)
					return
				}
				dequeuedFiles = append(dequeuedFiles, dequeuedFile.(*File))
			} else {
				break
			}
		}

		err = nm.encodeToIntermediateFile(dequeuedFiles, encoder)
		if err != nil {
			nm.logger.Error("Error encoding to reducer file: " + err.Error())
			nm.HandleProcessUpdate(proto, fragName, false, nil)
			return
		}

		keys = append(keys, dequeuedFiles[0].CurrKV.Key)

		for _, file := range dequeuedFiles {
			file.CurrKV = types.ReduceInputRecord{}
			err = file.Decoder.Decode(&file.CurrKV)
			if err != nil {
				if err == io.EOF {
					file.OpenFile.Close()
					err = nil
					continue
				} else {
					nm.logger.Error("Error decoding file: " + err.Error())
					nm.HandleProcessUpdate(proto, fragName, false, nil)
					return
				}
			}
			queue.Enqueue(file)
		}

	}

	err = nm.deleteIntermediateChunks(chunkFileNames)
	if err != nil {
		nm.logger.Error("Error deleting intermediate file: " + err.Error())
		nm.HandleProcessUpdate(proto, fragName, false, nil)
		return
	}

	//keys, err := nm.writeIntermediateFiles(output, fragPath)
	//if err != nil {
	//	nm.logger.Error("Error writing intermediate files: " + err.Error())
	//	nm.HandleProcessUpdate(proto, fragName, false, nil)
	//}
	output = nil

	nm.logger.Info("Map task completed")
	nm.HandleProcessUpdate(proto, fragName, true, keys)

}
*/

func (nm *NodeManager) deleteIntermediateChunks(chunkFileNames []string) (err error) {
	for _, chunkFileName := range chunkFileNames {
		err = os.Remove(chunkFileName)
		if err != nil {
			nm.logger.Error("Error deleting intermediate file: " + err.Error())
			return
		}
	}
	return
}

func (nm *NodeManager) encodeToIntermediateFile(dequeuedFiles []*File, encoder *gob.Encoder) (err error) {

	values := make([][]byte, 0)
	for _, file := range dequeuedFiles {
		values = append(values, file.CurrKV.Values...)
	}

	err = encoder.Encode(types.ReduceInputRecord{Key: dequeuedFiles[0].CurrKV.Key, Values: values})
	if err != nil {
		nm.logger.Error("Error encoding to file: " + err.Error())
		return err
	}

	return
}

type byKey []types.MapOutputRecord

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return bytes.Compare(a[i].Key, a[j].Key) < 0 }

func (nm *NodeManager) sortKeys(output []types.MapOutputRecord) (err error) {

	sort.Sort(byKey(output))

	nm.logger.Info("Sorted keys")

	return
}

func (nm *NodeManager) writeIntermediateFiles(intermediate []types.MapOutputRecord, fragPath string) (keys [][]byte, err error) {

	intermediateFile, err := os.Create(fragPath + "_MR_intermediate")
	if err != nil {
		nm.logger.Error("Error creating intermediate file: " + err.Error())

		return
	}
	defer intermediateFile.Close()
	encoder := gob.NewEncoder(intermediateFile)
	keys = make([][]byte, 0)

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && bytes.Equal(intermediate[j].Key, intermediate[i].Key) {
			j++
		}
		values := make([][]byte, 0)
		for k := i; k < j; k++ {
			keys = append(keys, intermediate[k].Key)
			values = append(values, intermediate[k].Value)
		}

		err = encoder.Encode(types.ReduceInputRecord{
			Key:    intermediate[i].Key,
			Values: values,
		})
		if err != nil {
			nm.logger.Error("Error encoding intermediate file: " + err.Error())
			return
		}
		i = j
	}

	nm.logger.Info("Intermediate file created")

	return
}
func (nm *NodeManager) processSetupRequest(response *resourceManagerProto3.SetupRequest, proto *resourceManagerProto3.ProtoHandler) {

	fileName := response.FileName
	job := response.Job

	/*

		Steps:
		- Check if file exists
			- If file exists, overwrite it (done)
			- If file does not exist, create it (done)
		- Write job to file (done)

	*/

	filePath := nm.dir + fileName + "_MR_job" + ".so"

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, fs.ModePerm)
	if err != nil {
		nm.logger.Error("Error creating file")
	}

	defer file.Close()

	file.Write(job)

	nm.logger.Info("Job written to file")

}
func (nm *NodeManager) HandleProcessUpdate(proto *resourceManagerProto3.ProtoHandler, fragName string, success bool, keys [][]byte) {

	proto.HandleProcessUpdate(fragName, success, keys)

}
