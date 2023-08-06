package resourceman_controller

import (
	"src/controller/storage_handler"
	messages "src/messages/resourceman_controller"
)

type RequestInterface interface {
	RequestType() string
}

type LayoutRequest struct {
	reqType  string
	FileName string
}

func (lr *LayoutRequest) RequestType() (request string) {
	return lr.reqType
}

func (p *ProtoHandler) fetchLayoutRequest(msg *messages.ResourceManagerRequest_LayoutRequest_) (req RequestInterface) {

	req = &LayoutRequest{
		reqType:  "layout",
		FileName: msg.LayoutRequest.FileName,
	}

	return

}

func (p *ProtoHandler) HandleLayoutResponse(fileMap map[string][]*storage_handler.Node) {

	p.logger.Info("Handling Get response to send.")
	var res *messages.ControllerResponse_FragLayoutResponse_

	if fileMap == nil {

		res = &messages.ControllerResponse_FragLayoutResponse_{
			FragLayoutResponse: &messages.ControllerResponse_FragLayoutResponse{
				ErrorCode: messages.ControllerResponse_FILE_NOT_FOUND,
			},
		}

	} else {

		res = &messages.ControllerResponse_FragLayoutResponse_{
			FragLayoutResponse: &messages.ControllerResponse_FragLayoutResponse{
				ErrorCode:         messages.ControllerResponse_NO_ERROR,
				TotalNumFragments: uint32(len(fileMap)),
				//repeated fragments
				FragmentLayout: []*messages.ControllerResponse_FragLayoutResponse_FragmentInfo{},
			},
		}

		for frag, nodes := range fileMap {
			fragInfo := &messages.ControllerResponse_FragLayoutResponse_FragmentInfo{
				FragmentId:     frag,
				StorageNodeIds: []*messages.ControllerResponse_FragLayoutResponse_StorageNodeInfo{},
			}

			for _, node := range nodes {
				nodeInfo := &messages.ControllerResponse_FragLayoutResponse_StorageNodeInfo{
					StorageNodeId: node.GetID(),
					Host:          node.GetAddress(),
					Port:          node.GetOpenPort(),
				}
				fragInfo.StorageNodeIds = append(fragInfo.StorageNodeIds, nodeInfo)
			}

			res.FragLayoutResponse.FragmentLayout = append(res.FragLayoutResponse.FragmentLayout, fragInfo)
		}

	}
	wrapper := &messages.ControllerResponse{
		ControllerMessage: res,
	}

	p.msgHandler.ServerResponseSend(wrapper)

}
