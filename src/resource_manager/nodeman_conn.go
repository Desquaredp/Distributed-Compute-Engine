package main

import (
	"net"
	messagesNodeManager "src/messages/resourceman_nodeman"
	proto3Controller "src/proto/resourceman_controller"
	proto3NodeManager "src/proto/resourceman_nodeman"
	"time"
)

func (rm *ResourceManager) DialNodeManager(id string) (msgHandler *messagesNodeManager.MessageHandler, err error) {
	conn, err := net.Dial("tcp", id)
	if err != nil {
		rm.logger.Sugar().Warn("There was error while dialing to the controller: \n", err.Error())
	}

	msgHandler = messagesNodeManager.NewMessageHandler(conn)

	return
}

func (rm *ResourceManager) DispatchTasks(loadBalancedMap map[string][]proto3Controller.StorageNodeInfo, job []byte, file string, responseChan chan proto3NodeManager.ResponseInterface) (dispatchedNodeSet map[proto3Controller.StorageNodeInfo]bool) {

	dispatchedSet := make(map[string]bool)
	dispatchedNodeSet = make(map[proto3Controller.StorageNodeInfo]bool)
	for fragmentId, nodeList := range loadBalancedMap {

		rm.logger.Info("Fragment: " + fragmentId)
		for _, node := range nodeList {

			rm.logger.Info("Sending request to Node Manager on node: " + node.Host)
			//TODO: REMOVE HARDCODED PORT
			msgHandler, err := rm.DialNodeManager(node.Host + ":" + "23099")
			if err != nil {
				rm.logger.Warn(err.Error())
				continue
			}

			go rm.handleNodeManager(msgHandler, responseChan)

			if err != nil {
				rm.logger.Warn(err.Error())
				continue
			} else {
				nodeManagerProtoInstance := proto3NodeManager.NewProtoHandler(msgHandler, rm.logger)
				if _, ok := dispatchedSet[node.NodeId]; ok {
					rm.HandleFragment(nodeManagerProtoInstance, fragmentId, file)
				} else {
					dispatchedNodeSet[node] = true
					dispatchedSet[node.NodeId] = true
					rm.HandleSetup(nodeManagerProtoInstance, file, job)
					//TODO: Tasks will fail if the setup is not done. Wait for the setup to complete.
					time.Sleep(5 * time.Second)
					rm.HandleFragment(nodeManagerProtoInstance, fragmentId, file)
				}

				break
			}

		}
	}

	return
}

func (rm *ResourceManager) handleNodeManager(msgHandler *messagesNodeManager.MessageHandler, responseChan chan proto3NodeManager.ResponseInterface) {
	proto := proto3NodeManager.NewProtoHandler(msgHandler, rm.logger)
	defer proto.MsgHandler().Close()

	for {

		wrapper, _ := proto.MsgHandler().NodeManagerResponseReceive()

		switch wrapper.Task.(type) {
		default:
			res := proto.HandleNodeManagerResponse(wrapper)
			rm.logger.Info("Received job response from the node manager.")

			switch res.(type) {

			case *proto3NodeManager.Progress:
				response := res.(*proto3NodeManager.Progress)
				rm.logger.Info("Progress: " + response.Status)

				responseChan <- response
				return
			case *proto3NodeManager.ShuffleResponse:
				response := res.(*proto3NodeManager.ShuffleResponse)
				if response.Success {
					rm.logger.Info("Shuffle successful")
				} else {
					rm.logger.Info("Shuffle failed")
				}
				responseChan <- response
				return
			case *proto3NodeManager.ReduceComplete:
				response := res.(*proto3NodeManager.ReduceComplete)
				if response.Success {
					rm.logger.Info("Reduce successful")
				} else {
					rm.logger.Info("Reduce failed")
				}
				responseChan <- response

			}

		case nil:
			return

		}

	}

}
