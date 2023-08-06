package main

import (
	"bytes"
	"crypto/md5"
	"net"
	"os"
	nodeManagerMessages "src/messages/nodeman_nodeman"
	nodeManagerProto3 "src/proto/nodeman_nodeman"
)

func (nm *NodeManager) AcceptNodeManagerConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			nm.logger.Fatal("Error accepting connection")
		} else {
			nm.logger.Info("Node Manager connected")
			msgHandler := nodeManagerMessages.NewMessageHandler(conn)
			go nm.HandleNodeManagerConnection(msgHandler)

		}

	}
}

func (nm *NodeManager) HandleNodeManagerConnection(msgHandler *nodeManagerMessages.MessageHandler) {
	proto := nodeManagerProto3.NewProtoHandler(msgHandler, nm.logger)
	defer proto.MsgHandler().Close()
	for {
		wrapper, err := proto.MsgHandler().ServerResponseReceive()
		if err != nil {
			nm.logger.Info("Node Manager connection closed")
			return
		}
		switch wrapper.Request.(type) {
		default:
			req := proto.HandleNodeManagerRequest(wrapper)
			nm.logger.Info("Received request from the Node Manager.")
			switch req.(type) {
			case *nodeManagerProto3.ReducerFileRequest:
				request := req.(*nodeManagerProto3.ReducerFileRequest)
				nm.logger.Info("Processing reducer file request")
				nm.processReducerFileRequest(request, proto)

				return
			}

		case nil:
			return

		}
	}

}

func (nm *NodeManager) processReducerFileRequest(request *nodeManagerProto3.ReducerFileRequest, proto *nodeManagerProto3.ProtoHandler) {

	nm.logger.Info("Processing reducer file request")
	nm.logger.Info("Reducer file name: " + request.FileName)

	checkSum := request.FileChecksum
	fileData := request.FileData

	md5CheckSum := md5.Sum(fileData)

	if bytes.Compare(md5CheckSum[:], checkSum[:]) == 0 {
		nm.logger.Info("File checksum matches")
	} else {
		nm.logger.Info("File checksum does not match")
		proto.SendReducerFileResponse(false)
	}

	file, err := os.Create(nm.dir + request.FileName)
	if err != nil {
		nm.logger.Error("Error creating file")
		proto.SendReducerFileResponse(false)
		return
	}

	defer file.Close()

	_, err = file.Write(fileData)
	if err != nil {
		nm.logger.Error("Error writing to file")
		proto.SendReducerFileResponse(false)
		return
	}

	nm.logger.Info("File written successfully")
	proto.SendReducerFileResponse(true)

}
