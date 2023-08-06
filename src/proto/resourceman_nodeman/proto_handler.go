package resourceman_nodeman

import (
	"go.uber.org/zap"
	messages "src/messages/resourceman_nodeman"
)

type ProtoHandler struct {
	msgHandler *messages.MessageHandler
	logger     *zap.Logger
	nodeId     string
}

func NewProtoHandler(msgHandler *messages.MessageHandler, logger *zap.Logger) *ProtoHandler {
	newProtoHandler := &ProtoHandler{
		msgHandler: msgHandler,
		logger:     logger,
	}
	return newProtoHandler
}

func (p *ProtoHandler) MsgHandler() *messages.MessageHandler {
	return p.msgHandler
}

func (p *ProtoHandler) HandleNodeManagerResponse(wrapper *messages.NodeManagerTaskResponse) (res ResponseInterface) {

	switch msg := wrapper.Task.(type) {

	case *messages.NodeManagerTaskResponse_Process_:
		p.logger.Info("Received progress response from the node manager.")
		res = p.fetchProcessResponse(msg)
		return

	case *messages.NodeManagerTaskResponse_ShuffleComplete_:
		p.logger.Info("Received shuffle complete response from the node manager.")
		res = p.fetchShuffleCompleteResponse(msg)
		return

	case *messages.NodeManagerTaskResponse_ReduceComplete_:
		p.logger.Info("Received reduce complete response from the node manager.")
		res = p.fetchReduceCompleteResponse(msg)
		return
	}

	return
}

func (p *ProtoHandler) HandleResourceManagerRequest(wrapper *messages.ResourceManagerTaskRequest) (req RequestInterface) {

	switch msg := wrapper.Task.(type) {

	case *messages.ResourceManagerTaskRequest_Setup_:
		p.logger.Info("Received Setup request from Resource Manager")
		req = p.fetchSetupRequest(msg)
		return

	case *messages.ResourceManagerTaskRequest_Fragment_:
		p.logger.Info("Received fragment request from Resource manager")
		req = p.fetchFragmentRequest(msg)
		return

	case *messages.ResourceManagerTaskRequest_KeyDistribution_:
		p.logger.Info("Received key distribution request from Resource manager")
		req = p.fetchKeyDistributionRequest(msg)
		return

	case *messages.ResourceManagerTaskRequest_ReduceTask_:
		p.logger.Info("Received reduce task request from Resource manager")
		req = p.fetchReduceTaskRequest(msg)
		return

	}

	return

}

func (p *ProtoHandler) HandleProcessUpdate(name string, success bool, keys [][]byte) {

	var req *messages.NodeManagerTaskResponse_Process_

	if success {
		p.logger.Info("Sending process update to Resource Manager")

		req = &messages.NodeManagerTaskResponse_Process_{
			Process: &messages.NodeManagerTaskResponse_Process{
				Status:       messages.NodeManagerTaskResponse_SUCCESS,
				FragmentName: name,
				Keys:         keys,
			},
		}

	} else {
		p.logger.Info("Sending process update to Resource Manager")

		req = &messages.NodeManagerTaskResponse_Process_{
			Process: &messages.NodeManagerTaskResponse_Process{
				Status:       messages.NodeManagerTaskResponse_FAILURE,
				FragmentName: name,
			},
		}
	}

	wrapper := &messages.NodeManagerTaskResponse{
		Task: req,
	}

	p.MsgHandler().NodeManagerResponseSend(wrapper)

}

func (p *ProtoHandler) HandleShuffleComplete(success bool) {

	var req *messages.NodeManagerTaskResponse_ShuffleComplete

	if success {
		p.logger.Info("Sending shuffle complete to Resource Manager")

		req = &messages.NodeManagerTaskResponse_ShuffleComplete{
			Status: messages.NodeManagerTaskResponse_SUCCESS,
		}

	} else {
		p.logger.Info("Sending shuffle complete to Resource Manager")

		req = &messages.NodeManagerTaskResponse_ShuffleComplete{
			Status: messages.NodeManagerTaskResponse_FAILURE,
		}
	}

	wrapper := &messages.NodeManagerTaskResponse{
		Task: &messages.NodeManagerTaskResponse_ShuffleComplete_{ShuffleComplete: req},
	}

	p.MsgHandler().NodeManagerResponseSend(wrapper)

}

func (p *ProtoHandler) HandleReduceRequest(inputFile string, outputFile string, reducers []int32) {

	req := &messages.ResourceManagerTaskRequest_ReduceTask_{
		ReduceTask: &messages.ResourceManagerTaskRequest_ReduceTask{
			InputFileName:  inputFile,
			OutputFileName: outputFile,
			TaskID:         reducers,
		},
	}

	wrapper := &messages.ResourceManagerTaskRequest{
		Task: req,
	}
	p.MsgHandler().ResourceManagerRequestSend(wrapper)
}

func (p *ProtoHandler) HandleReduceTaskUpdate(id int32, b bool) {

	var req *messages.NodeManagerTaskResponse_ReduceComplete_

	if b {
		p.logger.Info("Sending reduce task update to Resource Manager")

		req = &messages.NodeManagerTaskResponse_ReduceComplete_{
			ReduceComplete: &messages.NodeManagerTaskResponse_ReduceComplete{
				Status: messages.NodeManagerTaskResponse_SUCCESS,
				TaskID: id,
			},
		}

	} else {
		p.logger.Info("Sending reduce task update to Resource Manager")

		req = &messages.NodeManagerTaskResponse_ReduceComplete_{
			ReduceComplete: &messages.NodeManagerTaskResponse_ReduceComplete{
				Status: messages.NodeManagerTaskResponse_FAILURE,
				TaskID: id,
			},
		}
	}

	wrapper := &messages.NodeManagerTaskResponse{
		Task: req,
	}

	p.MsgHandler().NodeManagerResponseSend(wrapper)

}
