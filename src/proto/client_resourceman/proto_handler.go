package client_resourceman

import (
	"go.uber.org/zap"
	messages "src/messages/client_resourceman"
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

func (p *ProtoHandler) HandleResourceManagerResponse(wrapper *messages.ResourceManagerMessage) (res ResponseInterface) {

	switch msg := wrapper.Message.(type) {
	//switch _ := wrapper.Message.(type) {

	case *messages.ResourceManagerMessage_JobResponse_:
		p.logger.Info("Received job response from the ResourceManager.")
		res = p.fetchJobResponse(msg)
		return

	case *messages.ResourceManagerMessage_ProgressUpdate_:
		p.logger.Info("Received progress response from the ResourceManager.")
		res = p.fetchProgressResponse(msg)
		return
	}

	return
}

func (p *ProtoHandler) HandleClientRequest(wrapper *messages.ClientJobRequest) (req RequestInterface) {

	switch msg := wrapper.ClientMessage.(type) {

	case *messages.ClientJobRequest_Jobs_:
		p.logger.Info("Received job request from the client.")
		req = p.fetchJobRequest(msg)
		return

	}

	return

}

func (p *ProtoHandler) HandleProgressUpdate(count int, responses int, phase string, success bool) {

	var req *messages.ResourceManagerMessage_ProgressUpdate_

	if phase == "map" {
		taskEnum := messages.ResourceManagerMessage_MAP

		req = &messages.ResourceManagerMessage_ProgressUpdate_{
			ProgressUpdate: &messages.ResourceManagerMessage_ProgressUpdate{
				UpdateType: messages.ResourceManagerMessage_JOB_ONGOING,
				TaskProgress: &messages.ResourceManagerMessage_TaskProgress{
					TaskType:   taskEnum,
					TaskNumber: int32(count),
					TotalTasks: int32(responses),
				},
			},
		}

	} else if phase == "reduce" {
		taskEnum := messages.ResourceManagerMessage_REDUCE

		req = &messages.ResourceManagerMessage_ProgressUpdate_{
			ProgressUpdate: &messages.ResourceManagerMessage_ProgressUpdate{
				UpdateType: messages.ResourceManagerMessage_JOB_ONGOING,
				TaskProgress: &messages.ResourceManagerMessage_TaskProgress{
					TaskType:   taskEnum,
					TaskNumber: int32(count),
					TotalTasks: int32(responses),
				},
			},
		}
	} else if phase == "complete" {

		req = &messages.ResourceManagerMessage_ProgressUpdate_{
			ProgressUpdate: &messages.ResourceManagerMessage_ProgressUpdate{
				UpdateType: messages.ResourceManagerMessage_JOB_COMPLETED,
			},
		}

	}

	wrapper := &messages.ResourceManagerMessage{
		Message: req,
	}
	p.logger.Info("Sending progress update to the ResourceManager.")
	p.msgHandler.ResourceManagerResponseSend(wrapper)

}
