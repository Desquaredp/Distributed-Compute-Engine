package resourceman_nodeman

import messages "src/messages/resourceman_nodeman"

type RequestInterface interface {
	GetType() string
}

type SetupRequest struct {
	ReqType  string
	FileName string
	Job      []byte
}

func (sr *SetupRequest) GetType() string {
	return sr.ReqType
}

type FragmentRequest struct {
	ReqType  string
	Fragment string
	FileName string
}

func (fr *FragmentRequest) GetType() string {
	return fr.ReqType
}

type StorageNode struct {
	NodeID        string
	Host          string
	Port          string
	ReduceTaskIds []int32
}

type ReducerKeyDistribution struct {
	ReducerID int32
	Keys      [][]byte
}

type KeyDistributionRequest struct {
	ReqType                string
	StorageNode            []*StorageNode
	ReducerKeyDistribution []*ReducerKeyDistribution
}

func (kdr *KeyDistributionRequest) GetType() string {
	return kdr.ReqType
}

func (p *ProtoHandler) fetchSetupRequest(msg *messages.ResourceManagerTaskRequest_Setup_) (req RequestInterface) {
	req = &SetupRequest{
		ReqType:  "setup",
		FileName: msg.Setup.FileName,
		Job:      msg.Setup.Job,
	}

	return
}

func (p *ProtoHandler) fetchFragmentRequest(msg *messages.ResourceManagerTaskRequest_Fragment_) (req RequestInterface) {

	req = &FragmentRequest{
		ReqType:  "fragment",
		Fragment: msg.Fragment.FragmentName,
		FileName: msg.Fragment.FileName,
	}

	return

}

func (p *ProtoHandler) fetchKeyDistributionRequest(msg *messages.ResourceManagerTaskRequest_KeyDistribution_) (req RequestInterface) {

	req = &KeyDistributionRequest{
		ReqType:                "key_distribution",
		StorageNode:            make([]*StorageNode, 0),
		ReducerKeyDistribution: make([]*ReducerKeyDistribution, 0),
	}

	for _, node := range msg.KeyDistribution.StorageNode {
		req.(*KeyDistributionRequest).StorageNode = append(req.(*KeyDistributionRequest).StorageNode, &StorageNode{
			NodeID:        node.NodeID,
			Host:          node.Host,
			Port:          node.Port,
			ReduceTaskIds: node.ReduceTaskIDs,
		})
	}

	for _, reducer := range msg.KeyDistribution.ReducerKeyDistribution {
		req.(*KeyDistributionRequest).ReducerKeyDistribution = append(req.(*KeyDistributionRequest).ReducerKeyDistribution, &ReducerKeyDistribution{
			ReducerID: reducer.ReducerTaskID,
			Keys:      reducer.Keys,
		})
	}
	return
}

func (p *ProtoHandler) HandleSetupRequest(id string, job []byte) {

	req := &messages.ResourceManagerTaskRequest_Setup_{
		Setup: &messages.ResourceManagerTaskRequest_Setup{
			FileName: id,
			Job:      job,
		},
	}

	wrapper := &messages.ResourceManagerTaskRequest{
		Task: req,
	}

	p.MsgHandler().ResourceManagerRequestSend(wrapper)

}

func (p *ProtoHandler) HandleFragmentRequest(id string, name string) {

	req := &messages.ResourceManagerTaskRequest_Fragment_{
		Fragment: &messages.ResourceManagerTaskRequest_Fragment{
			FragmentName: id,
			FileName:     name,
		},
	}

	wrapper := &messages.ResourceManagerTaskRequest{
		Task: req,
	}

	p.MsgHandler().ResourceManagerRequestSend(wrapper)

}

type ReduceTaskRequest struct {
	ReqType      string
	InputFileName  string
	OutputFileName string
	TaskID         []int32
}

func (rtr *ReduceTaskRequest) GetType() string {
	return rtr.ReqType
}
func (p *ProtoHandler) fetchReduceTaskRequest(msg *messages.ResourceManagerTaskRequest_ReduceTask_) (req RequestInterface) {

	req = &ReduceTaskRequest{
		InputFileName:  msg.ReduceTask.InputFileName,
		OutputFileName: msg.ReduceTask.OutputFileName,
		TaskID:         msg.ReduceTask.TaskID,
	}

	return

}
