package main

import (
	"errors"
	proto3Controller "src/proto/resourceman_controller"
)

func (rm *ResourceManager) handleController(proto *proto3Controller.ProtoHandler) (res proto3Controller.ResponseInterface, err error) {
	defer proto.MsgHandler().Close()

	for {

		wrapper, _ := proto.MsgHandler().ServerResponseReceive()

		switch wrapper.ControllerMessage.(type) {

		default:
			res = proto.HandleControllerResponse(wrapper)

			switch res.ResponseType() {

			case "layout":
				response := res.(*proto3Controller.LayoutResponse)
				err = rm.processLayoutResponse(response)
				if err != nil {
					return
				} else {
					return response, nil
				}
			}

		case nil:
			return nil, errors.New("nil response")

		}

	}

}

func (rm *ResourceManager) processLayoutResponse(res *proto3Controller.LayoutResponse) (err error) {
	rm.logger.Info("Received layout response from the controller")

	if res.ErrorCode != "" {
		//TODO: propagate the file not found error to the client
		err = errors.New("error")
	} else {
		rm.logger.Info("Found file.")

	}
	return
}
