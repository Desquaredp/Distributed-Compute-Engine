package resourceman_controller

import messages "src/messages/resourceman_controller"

type ResponseInterface interface {
	ResponseType() string
}

type LayoutResponse struct {
	responseType      string
	ErrorCode         string
	TotalNumFragments uint32
	FragmentLayout    []FragmentInfo
}

type FragmentInfo struct {
	FragmentId   string
	Size         int64
	StorageNodes []StorageNodeInfo
}

type StorageNodeInfo struct {
	NodeId string
	Host   string
	Port   string
}

func (lr *LayoutResponse) ResponseType() string {
	return lr.responseType
}

func (p *ProtoHandler) fetchLayoutResponse(msg *messages.ControllerResponse_FragLayoutResponse_) (res ResponseInterface) {

	p.logger.Info("Received fragment layout response from the Controller.")
	p.logger.Sugar().Info("Status code: ", msg.FragLayoutResponse.ErrorCode.String())
	p.logger.Sugar().Info("Num of frags: ", msg.FragLayoutResponse.TotalNumFragments)

	if msg.FragLayoutResponse.ErrorCode != messages.ControllerResponse_NO_ERROR {
		res = &LayoutResponse{
			responseType: "layout",
			ErrorCode:    msg.FragLayoutResponse.ErrorCode.String(),
		}
		return
	}

	res = &LayoutResponse{
		responseType:      "layout",
		ErrorCode:         "",
		TotalNumFragments: msg.FragLayoutResponse.TotalNumFragments,
		FragmentLayout:    make([]FragmentInfo, 0),
	}

	i := 0
	for _, frag := range msg.FragLayoutResponse.FragmentLayout {

		p.logger.Sugar().Info("FragID: ", frag.FragmentId)
		p.logger.Sugar().Info("Frag Size: ", frag.Size)
		res.(*LayoutResponse).FragmentLayout = append(res.(*LayoutResponse).FragmentLayout, FragmentInfo{
			FragmentId:   frag.FragmentId,
			Size:         frag.Size,
			StorageNodes: make([]StorageNodeInfo, 0),
		})

		for _, node := range frag.StorageNodeIds {

			res.(*LayoutResponse).FragmentLayout[i].StorageNodes = append(res.(*LayoutResponse).FragmentLayout[i].StorageNodes, StorageNodeInfo{
				NodeId: node.StorageNodeId,
				Host:   node.Host,
				Port:   node.Port,
			})
		}
		i++
	}

	return

}

func (p *ProtoHandler) HandleLayoutRequest(file string) {

	req := &messages.ResourceManagerRequest{
		Message: &messages.ResourceManagerRequest_LayoutRequest_{

			LayoutRequest: &messages.ResourceManagerRequest_LayoutRequest{
				FileName: file,
			},
		},
	}
	p.logger.Info("Sending layout request to the controller")

	p.msgHandler.ClientRequestSend(req)

}
