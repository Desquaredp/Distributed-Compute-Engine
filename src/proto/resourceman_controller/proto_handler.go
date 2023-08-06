package resourceman_controller

import (
	"go.uber.org/zap"
	messages "src/messages/resourceman_controller"
)

type ProtoHandler struct {
	msgHandler *messages.MessageHandler
	logger     *zap.Logger
}

func NewProtoHandler(msgHandler *messages.MessageHandler, logger *zap.Logger) *ProtoHandler {
	newProtoHandler := &ProtoHandler{
		msgHandler: msgHandler,
		logger:     logger,
	}
	return newProtoHandler
}

func (p *ProtoHandler) HandleControllerResponse(wrapper *messages.ControllerResponse) (res ResponseInterface) {

	switch msg := wrapper.ControllerMessage.(type) {

	case *messages.ControllerResponse_FragLayoutResponse_:
		p.logger.Info("Received layout response from the Controller.")
		res = p.fetchLayoutResponse(msg)
		return

	}
	return
}

func (p *ProtoHandler) HandleResourceManagerRequest(wrapper *messages.ResourceManagerRequest) (req RequestInterface) {

	switch msg := wrapper.Message.(type) {

	case *messages.ResourceManagerRequest_LayoutRequest_:
		p.logger.Info("Received layout request from the resource manager")
		req = p.fetchLayoutRequest(msg)
		return

	}
	return
}

func (p *ProtoHandler) MsgHandler() *messages.MessageHandler {
	return p.msgHandler
}
