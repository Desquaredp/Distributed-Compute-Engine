package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	FileHandler "src/file"
	//	FileHandler "src/file"
)

func main() {

	inputType, err := parseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := appendLogger()
	defer logger.Sync()

	switch inputType.(type) {
	case *inputDFSYaml:

		fmt.Println("DFS job")
		dfsInput := inputType.(*inputDFSYaml)

		addr := dfsInput.Controller.Host + ":" + dfsInput.Controller.Port
		client := NewClient(addr, logger)

		client.Dial()

		fileName := dfsInput.InputFile
		fileHandler := FileHandler.NewFileHandler(fileName)
		fileHandler.SetDir(dfsInput.FileDir)
		fileHandler.CalcFileSize()

		client.HandlePUT(fileHandler, dfsInput.ChunkSize)

		client.HandleConnection()

	case *inputMRYaml:
		fmt.Println("MapReduce job")
		mrInput := inputType.(*inputMRYaml)

		addr := mrInput.ResourceManager.Host + ":" + mrInput.ResourceManager.Port
		client := NewClient(addr, logger)
		client.DialResourceManager()

		jobs := mrInput.Jobs
		inputFiles := make([]string, len(jobs))
		outputFile := make([]string, len(jobs))
		reducerCount := make([]int32, len(jobs))
		pluginPath := make([]string, len(jobs))

		for i, job := range jobs {

			fmt.Println("Job: ", job)

			inputFiles[i] = job.InputFile
			outputFile[i] = job.OutputFile
			reducerCount[i] = job.ReducerCount
			pluginPath[i] = job.JobBinary

		}

		pluginBinary := make([][]byte, 0)
		for _, plugin := range pluginPath {
			pluginBinary = append(pluginBinary, client.PopulateJobBinary(plugin))
		}

		client.HandleMRJob(inputFiles, outputFile, reducerCount, pluginBinary)
		client.HandleResourceManagerConnection()

	case nil:
		fmt.Println("No input type specified")
		return
	}

	/*
		operation := os.Args[1]

		if operation == "REST" {
			host := os.Args[2]

			client := NewClient(host, logger)
			client.Dial()

			fileName := os.Args[3]

			REST := os.Args[4]

			if REST == "PUT" {
				fileHandler := FileHandler.NewFileHandler(fileName)
				fileHandler.ExtractFileName()
				fileHandler.CalcFileSize()
				client.HandlePUT(fileHandler, 0)
				client.HandleConnection()
			} else if REST == "GET" {
				fileHandler := FileHandler.NewFileHandler(fileName)
				fileN := fileHandler.ExtractFileName()
				client.HandleGET(fileN)
				client.HandleConnection()
			}
		} else if operation == "mr" {

			host := os.Args[2]

			client := NewClient(host, logger)
			client.DialResourceManager()
			inputFile := os.Args[3]
			outputFile := os.Args[4]
			reducerCount := os.Args[5]
			reducerCountInt, _ := strconv.Atoi(reducerCount)
			jobPluginPath := os.Args[6]

			jobBinary := client.PopulateJobBinary(jobPluginPath)
			client.HandleMRJob(inputFile, outputFile, int32(reducerCountInt), jobBinary)
			client.HandleResourceManagerConnection()

		}
	*/

	select {}

}

type InputInterface interface {
	Type() string
}
type Address struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Job struct {
	InputFile    string `yaml:"input_file"`
	OutputFile   string `yaml:"output_file"`
	ReducerCount int32  `yaml:"reducer_count"`
	JobBinary    string `yaml:"job_binary"`
}

type inputMRYaml struct {
	ResourceManager Address `yaml:"resource_manager"`
	Jobs            []Job   `yaml:"jobs"`
}

func (i *inputMRYaml) Type() string {
	return "mr"
}

type inputDFSYaml struct {
	Controller Address `yaml:"controller"`
	InputFile  string  `yaml:"input_file"`
	FileDir    string  `yaml:"file_dir"`
	ChunkSize  int64   `yaml:"chunk_size"`
}

func (i *inputDFSYaml) Type() string {
	return "dfs"
}

func parseArgs(args []string) (inputType InputInterface, err error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments:\n use -h for help")
	}

	flag := args[1]

	switch flag {
	case "-h":

		fmt.Println("For help:")
		fmt.Println("./client -h")

		fmt.Println("To load config file:")
		fmt.Println("./client --load-config <dfs or mr> <config file>")

		fmt.Println("To populate config file with template values:")
		fmt.Println("./client --populate-config <dfs or mr> <config file>")

		os.Exit(0)

	case "--load-config":

		if len(args) < 4 {
			fmt.Println("not enough arguments")

			err = fmt.Errorf("not enough arguments:\n use --load-config <dfs or mr> <config file>")
			return
		}

		configType := args[2]
		configFile := args[3]

		if configType == "dfs" {
			//TODO: populate dfs config

			var readData inputDFSYaml
			readFile, err := os.Open(configFile)
			if err != nil {
				fmt.Printf("Error opening file: %s", err)
				return inputType, err
			}

			defer readFile.Close()
			decoder := yaml.NewDecoder(readFile)
			err = decoder.Decode(&readData)
			if err != nil {
				fmt.Printf("Error decoding YAML: %s", err)
				return inputType, err
			}

			fmt.Printf("Read data: %#v\n", readData)

			inputType = &readData
			return inputType, err

		} else if configType == "mr" {

			var readData inputMRYaml
			readFile, err := os.Open(configFile)
			if err != nil {
				fmt.Printf("Error opening file: %s", err)
				return inputType, err
			}
			defer readFile.Close()
			decoder := yaml.NewDecoder(readFile)
			err = decoder.Decode(&readData)
			if err != nil {
				fmt.Printf("Error decoding YAML: %s", err)
				return inputType, err
			}

			fmt.Printf("Read data: %#v\n", readData)
			inputType = &readData
			return inputType, err

		} else {
			fmt.Println("invalid config type")
			err = fmt.Errorf("invalid config type:\n use --load-config <dfs or mr> <config file>")
			return
		}

	case "--populate-config":
		if len(args) < 4 {
			fmt.Println("not enough arguments")

			err = fmt.Errorf("not enough arguments:\n use --populate-config <dfs or mr> <config file>")
			return
		}
		configType := args[2]
		configFile := args[3]

		if configType == "dfs" {
			//TODO: populate dfs config
			fmt.Println("Note: Only PUT operation is supported for now.")

			data := inputDFSYaml{
				Controller: Address{
					Host: "localhost",
					Port: "8080",
				},
				InputFile: "inputFile",
				FileDir:   "/path/to/file/",
				ChunkSize: 128000000,
			}

			file, err := os.Create(configFile)
			if err != nil {
				fmt.Printf("Error creating file: %s", err)
				return inputType, err
			}
			defer file.Close()
			encoder := yaml.NewEncoder(file)
			err = encoder.Encode(data)
			if err != nil {
				fmt.Printf("Error encoding YAML: %s", err)
				return inputType, err
			}

			fmt.Println("Config file populated successfully. Please edit the config file and provide the correct values.")
			os.Exit(0)

		} else if configType == "mr" {

			data := inputMRYaml{
				ResourceManager: Address{
					Host: "localhost",
					Port: "8080",
				},
				Jobs: []Job{
					{
						InputFile:    "inputFile",
						OutputFile:   "outputFile",
						ReducerCount: 5,
						JobBinary:    "/path/to/myjob",
					},
					{
						InputFile:    "input2",
						OutputFile:   "output2",
						ReducerCount: 3,
						JobBinary:    "/path/to/myjob2",
					},
				},
			}

			file, err := os.Create(configFile)
			if err != nil {
				fmt.Printf("Error creating file: %s", err)
				return inputType, err
			}
			defer file.Close()
			encoder := yaml.NewEncoder(file)
			err = encoder.Encode(data)
			if err != nil {
				fmt.Printf("Error encoding YAML: %s", err)
				return inputType, err
			}

			fmt.Println("Successfully created and populated " + configFile + " with template values")
			fmt.Println()
			fmt.Println("NOTE:")
			fmt.Println()

			fmt.Println("In case of job chaining:")
			fmt.Println("Job chaining requires jobs to be in order of execution in the config file. ")
			fmt.Print("Output file of one job will be used as input file for the next job. Please leave the input file of every job empty except the first job.")
			fmt.Println("To add more jobs, repeat the job structure in the config file.")

			fmt.Println()
			fmt.Println("In case of single job:")
			fmt.Println("Remove the second job from the config file.")

			fmt.Println()
			fmt.Println("In case of multiple mutually exclusive jobs:")
			fmt.Println("Provide the input file and output file for each job, and add more jobs by repeating the job structure in the config file, if necessary.")
			os.Exit(0)
			return inputType, err

		} else {

			err = fmt.Errorf("invalid config type:\n use --populate-config <dfs or mr> <config file>")
			return
		}

	}

	return
}

func appendLogger() *zap.Logger {

	file, err := os.OpenFile("logfileC.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Create a logger that writes to the file
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter := zapcore.AddSync(file)
	fileLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, fileLevel)

	// Create a logger that writes to the console
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriter := zapcore.Lock(os.Stdout)
	consoleLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, consoleLevel)

	// Create a final logger that writes to both the file and console
	logger := zap.New(zapcore.NewTee(fileCore, consoleCore))
	return logger
}
