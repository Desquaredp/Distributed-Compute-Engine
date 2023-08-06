package nodeman_nodeman

import (
	"go.uber.org/zap"
	messages "src/messages/nodeman_nodeman"
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

type RequestInterface interface {
	GetType() string
}

type ReducerFileRequest struct {
	ReqType      string
	FileName     string
	FileData     []byte
	FileChecksum []byte
}

func (rfr *ReducerFileRequest) GetType() string {
	return rfr.ReqType
}

func (p *ProtoHandler) HandleNodeManagerRequest(wrapper *messages.NodeManagerRequest) (req RequestInterface) {

	switch msg := wrapper.Request.(type) {

	case *messages.NodeManagerRequest_ReducerFile_:
		p.logger.Info("Received reducer file request from the node manager.")
		req = p.fetchReducerFileRequest(msg)
		return

	}

	return
}

type ResponseInterface interface {
	GetType() string
}

type ReducerFileResponse struct {
	ResType string
	Verdict bool
}

func (rfr *ReducerFileResponse) GetType() string {
	return rfr.ResType
}

func (p *ProtoHandler) HandleNodeManagerResponse(wrapper *messages.NodeManagerRequest) (res ResponseInterface) {

	switch msg := wrapper.Request.(type) {

	case *messages.NodeManagerRequest_ReducerFileResponse_:
		p.logger.Info("Received progress response from the node manager.")
		res = p.fetchProcessResponse(msg)
		return
	}

	return
}

func (p *ProtoHandler) fetchReducerFileRequest(msg *messages.NodeManagerRequest_ReducerFile_) (req RequestInterface) {

	req = &ReducerFileRequest{
		ReqType:      "reducer_file",
		FileName:     msg.ReducerFile.FileName,
		FileData:     msg.ReducerFile.FileData,
		FileChecksum: msg.ReducerFile.Checksum,
	}

	return
}

func (p *ProtoHandler) fetchProcessResponse(msg *messages.NodeManagerRequest_ReducerFileResponse_) (res ResponseInterface) {

	var verdict bool
	if msg.ReducerFileResponse.GetStatus() == messages.NodeManagerRequest_SUCCESS {
		verdict = true
	} else {
		verdict = false
	}

	res = &ReducerFileResponse{
		ResType: "reducer_file_response",
		Verdict: verdict,
	}

	return
}

func (p *ProtoHandler) SendReducerFileResponse(b bool) {

	var res *messages.NodeManagerRequest
	if b {
		p.logger.Info("Sending reducer file response to the node manager")
		res = &messages.NodeManagerRequest{
			Request: &messages.NodeManagerRequest_ReducerFileResponse_{
				ReducerFileResponse: &messages.NodeManagerRequest_ReducerFileResponse{
					Status: messages.NodeManagerRequest_SUCCESS,
				},
			},
		}

	} else {
		p.logger.Error("Sending reducer file response to the node manager")
		res = &messages.NodeManagerRequest{
			Request: &messages.NodeManagerRequest_ReducerFileResponse_{
				ReducerFileResponse: &messages.NodeManagerRequest_ReducerFileResponse{
					Status: messages.NodeManagerRequest_FAILURE,
				},
			},
		}
	}

	p.msgHandler.ClientRequestSend(res)

}

func (p *ProtoHandler) SendReducerFileTransfer(file string, data []byte, checksum [16]byte) {

	p.logger.Info("Sending reducer file transfer to the node manager")
	res := &messages.NodeManagerRequest{
		Request: &messages.NodeManagerRequest_ReducerFile_{
			ReducerFile: &messages.NodeManagerRequest_ReducerFile{
				FileName: file,
				FileData: data,
				Checksum: checksum[:],
			},
		},
	}

	p.msgHandler.ServerResponseSend(res)

}
