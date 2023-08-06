package resourceman_nodeman

import (
	messages "src/messages/resourceman_nodeman"
	proto3Controller "src/proto/resourceman_controller"
)

type ResponseInterface interface {
	GetType() string
}

type ShuffleResponse struct {
	ResType string
	Success bool
}

func (sr *ShuffleResponse) GetType() string {
	return sr.ResType
}

type Progress struct {
	ResType      string
	Status       string
	FragmentName string
	Keys         [][]byte
}

func (p *Progress) GetType() string {
	return p.ResType
}

type ReduceComplete struct {
	ResType   string
	ReducerID int32
	Success   bool
}

func (rc *ReduceComplete) GetType() string {
	return rc.ResType
}

func (p *ProtoHandler) fetchProcessResponse(msg *messages.NodeManagerTaskResponse_Process_) (res ResponseInterface) {

	if msg.Process.Status == messages.NodeManagerTaskResponse_SUCCESS {
		p.logger.Info("Received success response from the node manager.")
		res = &Progress{
			ResType:      "progress",
			Status:       "success",
			FragmentName: msg.Process.FragmentName,
			Keys:         msg.Process.Keys,
		}
	} else {
		p.logger.Info("Received failure response from the node manager.")
		res = &Progress{
			ResType:      "progress",
			Status:       "failure",
			FragmentName: msg.Process.FragmentName,
		}
	}

	return
}

func (p *ProtoHandler) DispatchKeyDistribution(nodes map[proto3Controller.StorageNodeInfo][]int32, reducerMap map[int][][]byte) {

	p.logger.Info("Sending key distribution to Resource Manager")

	req := &messages.ResourceManagerTaskRequest_KeyDistribution_{
		KeyDistribution: &messages.ResourceManagerTaskRequest_KeyDistribution{
			StorageNode:            make([]*messages.ResourceManagerTaskRequest_StorageNode, 0),
			ReducerKeyDistribution: make([]*messages.ResourceManagerTaskRequest_ReducerKeyDistribution, 0),
		},
	}

	for node, reduceTaskIds := range nodes {
		req.KeyDistribution.StorageNode = append(req.KeyDistribution.StorageNode, &messages.ResourceManagerTaskRequest_StorageNode{
			NodeID:        node.NodeId,
			Host:          node.Host,
			Port:          node.Port,
			ReduceTaskIDs: reduceTaskIds,
		})
	}

	for reducerId, keys := range reducerMap {
		req.KeyDistribution.ReducerKeyDistribution = append(req.KeyDistribution.ReducerKeyDistribution, &messages.ResourceManagerTaskRequest_ReducerKeyDistribution{
			ReducerTaskID: int32(reducerId),
			Keys:          keys,
		})
	}

	wrapper := &messages.ResourceManagerTaskRequest{
		Task: req,
	}

	p.MsgHandler().ResourceManagerRequestSend(wrapper)

}

func (p *ProtoHandler) fetchShuffleCompleteResponse(msg *messages.NodeManagerTaskResponse_ShuffleComplete_) (res ResponseInterface) {

	p.logger.Info("Received shuffle complete response from the node manager.")

	outcome := msg.ShuffleComplete.Status
	success := false

	if outcome == messages.NodeManagerTaskResponse_SUCCESS {
		p.logger.Info("Shuffle completed successfully.")
		success = true

	} else {
		p.logger.Info("Shuffle failed.")
	}
	return &ShuffleResponse{
		ResType: "shuffle",
		Success: success,
	}

}

func (p *ProtoHandler) fetchReduceCompleteResponse(msg *messages.NodeManagerTaskResponse_ReduceComplete_) (res ResponseInterface) {

	p.logger.Info("Received reduce complete response from the node manager.")

	outcome := msg.ReduceComplete.Status
	success := false

	if outcome == messages.NodeManagerTaskResponse_SUCCESS {
		p.logger.Info("Reduce completed successfully.")
		success = true

	} else {
		p.logger.Info("Reduce failed.")
	}
	return &ReduceComplete{
		ResType: "reduce",
		Success: success,
	}

}
