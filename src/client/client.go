package main

import (
	"go.uber.org/zap"
	"net"
	"src/file"
	messages "src/messages/controller_client"
	proto3Resourceman "src/proto/client_resourceman"
	proto3 "src/proto/controller_client"
)

type Client struct {
	serverPort       string
	file             *file.FileHandler
	logger           *zap.Logger
	conn             net.Conn
	proto            *proto3.ProtoHandler
	msgHandler       *messages.MessageHandler
	protoResourceman *proto3Resourceman.ProtoHandler
}

type File struct {
	name      string
	size      int64
	chunkSize int64
	checkSum  [32]byte
}

func NewClient(serverPort string, logger *zap.Logger) (client *Client) {

	client = &Client{
		serverPort: serverPort,
		logger:     logger,
	}

	return
}

func (c *Client) Disconnect() {
	c.conn.Close()
}
func (c *Client) Conn(conn net.Conn) {
	c.conn = conn
}

func (c *Client) MsgHandler(handler *messages.MessageHandler) {
	c.msgHandler = handler
}

func (c *Client) Proto(proto *proto3.ProtoHandler) {
	c.proto = proto
}

func (c *Client) Dial() (err error) {
	server := c.serverPort

	conn, err := net.Dial("tcp", server)

	if err != nil {
		c.logger.Fatal("There was an error connecting to the host.")
		panic("There was an error connecting to the host.")
		return
	}
	msgHandler := messages.NewMessageHandler(conn)
	c.MsgHandler(msgHandler)
	proto := proto3.NewProtoHandler(msgHandler, c.logger)
	c.Proto(proto)
	c.Conn(conn)

	return
}

func (c *Client) HandleConnection() (err error) {
	for {
		wrapper, _ := c.proto.MsgHandler().ControllerResponseReceive()

		switch wrapper.ControllerMessage.(type) {
		default:
			res := c.proto.HandleControllerResponse(wrapper)

			switch res.GetResType() {
			case "PlanResponse":
				if res.(*proto3.PlanResponse).StatusCode == "OK" {
					c.DispatchFile(res)
				}
			case "FragmentLayoutResponse":
				if res.(*proto3.FragLayoutResponse).StatusCode == "OK" {
					c.FetchFile(res)
				}

			}
		case nil:
			return

		}
	}
}

func (c *Client) HandlePUT(file *file.FileHandler, fragSize int64) (err error) {

	c.file = file
	c.proto.HandlePutRequest(file.ExtractFileName(), file.FileSize(), fragSize)
	return
}

func (c *Client) HandleGET(file string) {

	c.logger.Info("Handling GET request")
	c.proto.HandleGetRequest(file)

}
