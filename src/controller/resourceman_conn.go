package main

import (
	"go.uber.org/zap"
	"net"
	"src/controller/storage_handler"
	resourceManagerMessages "src/messages/resourceman_controller"
	resourceManagerProto3 "src/proto/resourceman_controller"
)

func acceptResourceManagerConnections(listener net.Listener, spokeHandler *storage_handler.StorageNodeHandler, logger *zap.Logger) {

	for {
		if conn, err := listener.Accept(); err == nil {

			logger.Info("Accepted a Resource Manager Connection")
			msgHandler := resourceManagerMessages.NewMessageHandler(conn)
			go handleResourceManager(msgHandler, spokeHandler, logger)

		} else {

			logger.Error(err.Error())
		}
	}

}

func handleResourceManager(msgHandler *resourceManagerMessages.MessageHandler, spokeHandler *storage_handler.StorageNodeHandler, logger *zap.Logger) {
	defer msgHandler.Close()
	proto := resourceManagerProto3.NewProtoHandler(msgHandler, logger)

	for {
		wrapper, _ := proto.MsgHandler().ClientRequestReceive()

		switch wrapper.Message.(type) {

		default:
			req := proto.HandleResourceManagerRequest(wrapper)

			if req == nil {
				//TODO: throw error or disconnect client
			}

			switch req.RequestType() {

			case "layout":
				request := req.(*resourceManagerProto3.LayoutRequest)

				logger.Info("Received layout request from the resource manager")
				logger.Sugar().Info(
					"Layout request for file: ",
					req.(*resourceManagerProto3.LayoutRequest).FileName,
				)

				logger.Info("Processing GET request")

				FileMap := spokeHandler.FindFiles(request.FileName, logger)
				if FileMap == nil {
					logger.Info("File doesn't exists.")
					proto.HandleLayoutResponse(nil)
				} else {
					logger.Info("File exists.")

					logger.Sugar().Info("FileMap length: ", len(FileMap))
					proto.HandleLayoutResponse(FileMap)
				}
			}
		case nil:
			logger.Info("Resource Manager Disconnected")
			return

		}
	}

}
