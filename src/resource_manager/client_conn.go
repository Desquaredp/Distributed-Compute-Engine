package main

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"os"
	"sort"
	clientMessages "src/messages/client_resourceman"
	clientProto3 "src/proto/client_resourceman"
	proto3Controller "src/proto/resourceman_controller"
	proto3NodeManager "src/proto/resourceman_nodeman"
	"strconv"
)

// acceptClientConnections accepts client connections and spawns a goroutine to handle each client.
func (rm *ResourceManager) acceptClientConnections(listener net.Listener) {
	rm.logger.Info("Accepting client connections.")

	for {
		if conn, err := listener.Accept(); err == nil {

			rm.logger.Info("Client connected.")
			msgHandler := clientMessages.NewMessageHandler(conn)
			go rm.handleClient(msgHandler)

		} else {
			rm.logger.Error(err.Error())
		}

	}

}

// handleClient handles client requests.
func (rm *ResourceManager) handleClient(msgHandler *clientMessages.MessageHandler) {

	defer msgHandler.Close()
	proto := clientProto3.NewProtoHandler(msgHandler, rm.logger)

	for {

		wrapper, _ := proto.MsgHandler().ClientRequestReceive()

		switch wrapper.ClientMessage.(type) {

		default:
			req := proto.HandleClientRequest(wrapper)
			rm.logger.Info("Received job request from the client.")

			switch req.(type) {
			case *clientProto3.JobRequest:
				request := req.(*clientProto3.JobRequest)

				if request.Jobs == nil {
					rm.logger.Warn("No jobs specified.")
					rm.logger.Warn("Closing client connection.")
					os.Exit(0)
					return
				}
				rm.logger.Info("Processing job request")
				rm.processJobRequest(request, proto)
				return
			}

		case nil:
			rm.logger.Info("Client connection closed")
			return
		}

	}

}

// processJobRequest processes job requests received from the client.
func (rm *ResourceManager) processJobRequest(req *clientProto3.JobRequest, proto *clientProto3.ProtoHandler) {

	clientJobRequest := req.Jobs
	var jobOutputFile string

	for _, jobInstance := range clientJobRequest {

		inputFile := jobInstance.InputFile

		if inputFile == "" {
			rm.logger.Info("Input file not specified. Possible job chaining detected.")
			rm.logger.Info(" Using the output file of the previous job as the input file of the current job.")
			inputFile = jobOutputFile
		}

		outputFile := jobInstance.OutputFile
		jobOutputFile = outputFile

		reducers := int(jobInstance.ReducerCount)
		job := jobInstance.JobBinary

		rm.logger.Info("Input file: " + inputFile)
		rm.logger.Info("Output file: " + outputFile)

		//TODO: Dry run the binary and check if it is valid?

		msgHandler := rm.DialController()
		controllerProtoInstance := proto3Controller.NewProtoHandler(msgHandler, rm.logger)
		controllerResponseChan := make(chan proto3Controller.ResponseInterface)
		rm.HandleLayout(controllerProtoInstance, inputFile)
		go func() {
			response, err := rm.handleController(controllerProtoInstance)
			if err != nil {
				rm.logger.Warn(err.Error())
				response = nil

			}
			controllerResponseChan <- response
			close(controllerResponseChan)

		}()

		response := <-controllerResponseChan

		if response == nil {
			rm.logger.Warn("Controller couldn't find the file.")
			proto.HandleJobResponse(errors.New("FILE_NOT_FOUND"))
			return
		} else {

			/*
				Steps:
					- Use the fragment layout to create a load distributor (done)
					- According to the load distribution, contact the node manager (done)
					- Send the task to the node manager (done)
			*/

			loadBalancedMap := rm.LoadDistributor(response.(*proto3Controller.LayoutResponse).FragmentLayout)

			rm.logger.Info("Load balanced map: ", zap.Any("loadBalancedMap", loadBalancedMap))

			responseChan := make(chan proto3NodeManager.ResponseInterface)

			// Set the total number of expected responses
			totalResponses := len(loadBalancedMap)

			// Start a goroutine to receive responses and update the client

			//Wait for this goroutine to finish before returning

			done := make(chan bool)
			keys := make([][]byte, 0, len(loadBalancedMap))
			go func() {
				count := 1
				for response := range responseChan {
					// ... process the response ...

					res := response.(*proto3NodeManager.Progress)

					if res.Status == "success" {
						rm.logger.Info("Fragment " + res.FragmentName + " processed successfully.")
						phase := "map"
						success := true
						if res.Keys != nil {
							keys = append(keys, res.Keys...)
						}

						rm.HandleProgressUpdate(proto, count, totalResponses, phase, success)
						count++
						if count == totalResponses+1 {
							rm.logger.Info("All fragments processed successfully.")
							done <- true
						}
					} else {
						rm.logger.Info("Fragment " + res.FragmentName + " failed.")
						phase := "map"
						success := false
						rm.HandleProgressUpdate(proto, count, totalResponses, phase, success)
						break

						//TODO: Handle failure. Retry? OR send the error to the client and kill?
					}

				}
				if count == totalResponses+1 {
					rm.logger.Info("All fragments processed successfully.")
					done <- true
				} else {
					rm.logger.Info("Some fragments failed.")
					done <- false
				}

			}()
			dispatchedNodesSet := rm.DispatchTasks(loadBalancedMap, job, inputFile, responseChan)
			//TODO: take the number of reducers from the client
			//reducers := len(dispatchedNodesSet)

			success := <-done
			if success {
				rm.logger.Info("All fragments processed successfully.")
				//TODO: move to the shuffle phase
				rm.logger.Info("Moving to the shuffle phase.")

				sort.Slice(keys, func(i, j int) bool {
					return false
				})
				reducerMap := rm.distributeKeys(keys, reducers)
				rm.logger.Info("Keys distributed to reducers.")
				assignedNodes := rm.assignReducers(reducers, dispatchedNodesSet)
				rm.logger.Info("Reducers assigned to nodes.")
				rm.logger.Info("Printing assigned nodes.")
				for node, reducers := range assignedNodes {
					rm.logger.Info("Node: " + node.Host + ":" + node.Port)
					rm.logger.Info("Reducers: ")
					for _, reducer := range reducers {
						rm.logger.Info(strconv.Itoa(int(reducer)))
					}
				}

				rm.logger.Info("Printing reducer map.")
				for reducer, keys := range reducerMap {
					rm.logger.Info("Reducer: " + strconv.Itoa(reducer))
					rm.logger.Info("Keys length: ")
					rm.logger.Info(strconv.Itoa(len(keys)))
				}

				rm.logger.Info("Dispatching tasks to reducers.")
				/*
					reducerResponseChan := make(chan *proto3NodeManager.Progress)
					reducerDone := make(chan bool)

				*/

				responseChan = make(chan proto3NodeManager.ResponseInterface)
				totalResponses = len(dispatchedNodesSet)
				rm.handleKMapDispatch(assignedNodes, reducerMap, dispatchedNodesSet, responseChan)

				go func() {

					count := 0
					for response := range responseChan {

						res := response.(*proto3NodeManager.ShuffleResponse)

						if res.Success {
							rm.logger.Info("Shuffle for node completed successfully.")
							count++
							if count == totalResponses {
								rm.logger.Info("All shuffles completed successfully.")
								done <- true
							}
						} else {
							rm.logger.Info("Shuffle for node failed.")
							done <- false
							break
						}

					}

				}()

				success = <-done
				if success {
					//TODO:move to the reduce phase
					rm.logger.Info("Moving to the reduce phase.")

					responseChan = make(chan proto3NodeManager.ResponseInterface)
					totalResponses = len(assignedNodes)

					rm.handleReduceDispatch(assignedNodes, inputFile, outputFile, responseChan)

					go func() {

						count := 0
						for response := range responseChan {

							res := response.(*proto3NodeManager.ReduceComplete)

							if res.Success {
								rm.logger.Info("Reduce task processed successfully.")
								phase := "reduce"

								count++
								rm.HandleProgressUpdate(proto, count, totalResponses, phase, success)
								if count == totalResponses {
									rm.logger.Info("All reducers processed successfully.")
									done <- true
								}
							} else {
								rm.logger.Info("Reduce task failed.")
								phase := "reduce"
								rm.HandleProgressUpdate(proto, count, totalResponses, phase, success)
								break

								//TODO: Handle failure. Retry? OR send the error to the client and kill?
							}

						}
						if count == totalResponses {
							rm.logger.Info("All fragments processed successfully.")
							done <- true
						} else {
							rm.logger.Info("Some reducers failed.")
							done <- false
						}

					}()

					success = <-done
					if success {
						rm.logger.Info("All reducers processed successfully.")
						rm.HandleProgressUpdate(proto, 0, 0, "complete", success)

					} else {
						rm.logger.Info("Some reducers failed.")
						proto.HandleJobResponse(errors.New("SERVER_ERROR"))
						return
					}
				} else {
					rm.logger.Info("Sending error to client.")
					proto.HandleJobResponse(errors.New("SERVER_ERROR"))
					return
				}

			} else {
				rm.logger.Info("Some fragments failed.")
				proto.HandleJobResponse(errors.New("SERVER_ERROR"))
				return
			}

		}
	}
	proto.HandleJobResponse(nil)

}

func (rm *ResourceManager) handleReduceDispatch(assignedNodes map[proto3Controller.StorageNodeInfo][]int32, inputFile string, outputFile string, responseChan chan proto3NodeManager.ResponseInterface) {

	for node, reducers := range assignedNodes {
		msgHandler, err := rm.DialNodeManager(node.Host + ":" + "23099")
		if err != nil {
			rm.logger.Warn("Couldn't connect to the node manager.")
		}
		nodeManagerInstance := proto3NodeManager.NewProtoHandler(msgHandler, rm.logger)
		go rm.handleNodeManager(msgHandler, responseChan)
		rm.logger.Info("Sending reduce request to node manager.")
		nodeManagerInstance.HandleReduceRequest(inputFile, outputFile, reducers)

	}
}
func (rm *ResourceManager) handleKMapDispatch(assignedNodes map[proto3Controller.StorageNodeInfo][]int32, reducerMap map[int][][]byte, set map[proto3Controller.StorageNodeInfo]bool, responseChan chan proto3NodeManager.ResponseInterface) {
	/*
				Steps:
					- Contact the node managers ()
			        - Send the info regarding the reducers and the keys to the node managers
			        - On nodeMan:
						- take the intermediate file
						- collect the keys based on the reducers and save them in their respective files
						- ship em to the respective reduce holder
					- On Reducer holder:
						- collect the files sent by the mappers
						- decide how to read.
		think about how many map files are generated


	*/

	for node := range set {
		msgHandler, err := rm.DialNodeManager(node.Host + ":" + "23099")
		if err != nil {
			rm.logger.Warn("Couldn't connect to the node manager.", zap.Error(err))
		}
		nodeManagerProtoInstance := proto3NodeManager.NewProtoHandler(msgHandler, rm.logger)
		go rm.handleNodeManager(msgHandler, responseChan)
		rm.logger.Info("Dispatching keys to node manager.")
		nodeManagerProtoInstance.DispatchKeyDistribution(assignedNodes, reducerMap)

	}

}

func (rm *ResourceManager) assignReducers(reducers int, set map[proto3Controller.StorageNodeInfo]bool) (assignedNodes map[proto3Controller.StorageNodeInfo][]int32) {

	assignedNodes = make(map[proto3Controller.StorageNodeInfo][]int32)
	if reducers <= len(set) {

		var i int32
		i = 0

		for node := range set {
			assignedNodes[node] = append(assignedNodes[node], i)
			i++
			if i >= int32(reducers) {
				break
			}
		}

	} else {
		dist := reducers / len(set)
		rem := reducers % len(set)
		var i int32
		i = 0
		for node := range set {
			for j := 0; j < dist; j++ {
				assignedNodes[node] = append(assignedNodes[node], i)
				i++
			}
			if rem > 0 {
				assignedNodes[node] = append(assignedNodes[node], i)
				i++
				rem--
			}

			if i > int32(reducers) {
				break
			}
		}

	}

	return
}

func (rm *ResourceManager) distributeKeys(keys [][]byte, reducers int) (reducerMap map[int][][]byte) {

	reducerMap = make(map[int][][]byte)
	for _, key := range keys {

		hash := md5.Sum(key)
		digest := binary.BigEndian.Uint64(hash[:])
		pos := int(digest % uint64(reducers))
		reducerMap[pos] = append(reducerMap[pos], key)
	}
	return
}

func (rm *ResourceManager) HandleProgressUpdate(proto *clientProto3.ProtoHandler, count int, responses int, phase string, success bool) {

	proto.HandleProgressUpdate(count, responses, phase, success)

}

func (rm *ResourceManager) HandleFragment(proto *proto3NodeManager.ProtoHandler, fragmentId string, fileName string) {

	proto.HandleFragmentRequest(fragmentId, fileName)

}

func (rm *ResourceManager) HandleSetup(proto *proto3NodeManager.ProtoHandler, fragmentId string, job []byte) {

	proto.HandleSetupRequest(fragmentId, job)

}

func (rm *ResourceManager) LoadDistributor(layout []proto3Controller.FragmentInfo) (distMap map[string][]proto3Controller.StorageNodeInfo) {

	distMap = make(map[string][]proto3Controller.StorageNodeInfo)
	for _, fragment := range layout {

		rm.logger.Info("Fragment: " + fragment.FragmentId)
		/*
			Steps:
				- shuffle the storage nodes(done)
				- Send the task to the node manager
		*/

		rand.Shuffle(len(fragment.StorageNodes), func(i, j int) {
			fragment.StorageNodes[i], fragment.StorageNodes[j] = fragment.StorageNodes[j], fragment.StorageNodes[i]
		})

		distMap[fragment.FragmentId] = fragment.StorageNodes

	}

	return
}

func (rm *ResourceManager) HandleLayout(proto *proto3Controller.ProtoHandler, inputFile string) {
	proto.HandleLayoutRequest(inputFile)
}
