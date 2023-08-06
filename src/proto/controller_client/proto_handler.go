package controller_client

import (
	"go.uber.org/zap"
	"src/controller/file_distributor"
	"src/controller/storage_handler"
	messages "src/messages/controller_client"
)

type ProtoHandler struct {
	msgHandler *messages.MessageHandler
	logger     *zap.Logger
	nodeId     string
}

func (p *ProtoHandler) MsgHandler() *messages.MessageHandler {
	return p.msgHandler
}

func NewProtoHandler(msgHandler *messages.MessageHandler, logger *zap.Logger) *ProtoHandler {
	newProtoHandler := &ProtoHandler{
		msgHandler: msgHandler,
		logger:     logger,
	}
	return newProtoHandler
}

func (p *ProtoHandler) HandleControllerResponse(wrapper *messages.ControllerMessage) (res ResponseInterface) {

	switch msg := wrapper.ControllerMessage.(type) {

	case *messages.ControllerMessage_PlanResponse_:
		p.logger.Info("Received plan response from the Controller.")
		res = p.fetchPlanResponse(msg)
		return
	case *messages.ControllerMessage_FragLayoutResponse_:
		res = p.fetchLayoutResponse(msg)
	case *messages.ControllerMessage_DeleteResponse_:
		p.fetchDeleteResponse(msg)
	case *messages.ControllerMessage_LsResponse_:
		p.fetchLsResponse(msg)

	}

	return
}

func (p *ProtoHandler) HandleClientRequest(wrapper *messages.ClientMessage) (req *Request) {

	switch msg := wrapper.ClientMessage.(type) {

	case *messages.ClientMessage_PutRequest_:
		req = p.fetchPutRequest(msg)

	case *messages.ClientMessage_GetRequest_:
		req = p.fetchGetRequest(msg)

	case *messages.ClientMessage_DeleteRequest_:
		p.fetchDeleteRequest(msg)

	case *messages.ClientMessage_LsRequest_:
		p.fetchLsRequest(msg)

	}

	return
}

func (p *ProtoHandler) HandlePlanResponse(fragMap map[*file_distributor.Fragment][]*storage_handler.Node, req *Request) {

	p.logger.Info("Sending plan response to the Controller.")

	res := &messages.ControllerMessage_PlanResponse_{
		PlanResponse: &messages.ControllerMessage_PlanResponse{
			StatusCode:        messages.ControllerMessage_OK,
			TotalNumFragments: uint32(len(fragMap)),
			//repeated fragments
			FragmentLayout: []*messages.ControllerMessage_PlanResponse_FragmentInfo{},
		},
	}

	for frag, nodes := range fragMap {
		fragInfo := &messages.ControllerMessage_PlanResponse_FragmentInfo{
			FragmentId:     frag.GetFragmentName(),
			Size:           frag.GetFragmentSize(),
			StorageNodeIds: []*messages.ControllerMessage_PlanResponse_StorageNodeInfo{},
		}

		for _, node := range nodes {
			nodeInfo := &messages.ControllerMessage_PlanResponse_StorageNodeInfo{
				StorageNodeId: node.GetID(),
				Host:          node.GetAddress(),
				Port:          node.GetOpenPort(),
			}
			fragInfo.StorageNodeIds = append(fragInfo.StorageNodeIds, nodeInfo)
		}

		res.PlanResponse.FragmentLayout = append(res.PlanResponse.FragmentLayout, fragInfo)
	}

	wrapper := &messages.ControllerMessage{
		ControllerMessage: res,
	}

	p.msgHandler.ControllerResponseSend(wrapper)

}

func (p *ProtoHandler) HandleGetRequest(file string) {

	p.logger.Info("Sending Get request to the Controller.")

	req := &messages.ClientMessage_GetRequest_{
		GetRequest: &messages.ClientMessage_GetRequest{
			RestOption: messages.ClientMessage_GET,
			FileName:   file,
		},
	}

	wrapper := &messages.ClientMessage{
		ClientMessage: req,
	}

	p.msgHandler.ClientRequestSend(wrapper)

}

func (p *ProtoHandler) HandleGetResponse(fileMap map[string][]*storage_handler.Node, req *Request) {

	p.logger.Info("Handling Get response to send.")
	var res *messages.ControllerMessage_FragLayoutResponse_

	if fileMap == nil {

		res = &messages.ControllerMessage_FragLayoutResponse_{
			FragLayoutResponse: &messages.ControllerMessage_FragLayoutResponse{
				StatusCode: messages.ControllerMessage_FILE_NOT_FOUND,
			},
		}

	} else {

		res = &messages.ControllerMessage_FragLayoutResponse_{
			FragLayoutResponse: &messages.ControllerMessage_FragLayoutResponse{
				StatusCode:        messages.ControllerMessage_OK,
				TotalNumFragments: uint32(len(fileMap)),
				//repeated fragments
				FragmentLayout: []*messages.ControllerMessage_FragLayoutResponse_FragmentInfo{},
			},
		}

		for frag, nodes := range fileMap {
			fragInfo := &messages.ControllerMessage_FragLayoutResponse_FragmentInfo{
				FragmentId:     frag,
				StorageNodeIds: []*messages.ControllerMessage_FragLayoutResponse_StorageNodeInfo{},
			}

			for _, node := range nodes {
				nodeInfo := &messages.ControllerMessage_FragLayoutResponse_StorageNodeInfo{
					StorageNodeId: node.GetID(),
					Host:          node.GetAddress(),
					Port:          node.GetOpenPort(),
				}
				fragInfo.StorageNodeIds = append(fragInfo.StorageNodeIds, nodeInfo)
			}

			res.FragLayoutResponse.FragmentLayout = append(res.FragLayoutResponse.FragmentLayout, fragInfo)
		}

	}
	wrapper := &messages.ControllerMessage{
		ControllerMessage: res,
	}

	p.msgHandler.ControllerResponseSend(wrapper)

}
