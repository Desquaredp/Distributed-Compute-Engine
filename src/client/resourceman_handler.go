package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	messagesResourceman "src/messages/client_resourceman"
	proto3Resourceman "src/proto/client_resourceman"
)

// SetProtoResourceman sets the proto handler for the resource manager.
func (c *Client) SetProtoResourceman(proto *proto3Resourceman.ProtoHandler) {
	c.protoResourceman = proto
}

// PopulateJobBinary populates the job binary.
func (c *Client) PopulateJobBinary(JobBinaryPath string) (jobBinary []byte) {

	file, err := os.Open(JobBinaryPath)
	if err != nil {
		c.logger.Fatal("There was an error opening the job binary.")
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		jobBinary = append(jobBinary, buffer[:n]...)

	}

	return
}

// HandleResourceManagerConnection  handles map reduce job comms with the resource manager
func (c *Client) HandleResourceManagerConnection() {
	defer c.protoResourceman.MsgHandler().Close()

	for {

		wrapper, _ := c.protoResourceman.MsgHandler().ResourceManagerResponseReceive()

		switch wrapper.Message.(type) {

		default:
			c.logger.Info("Received job response from the ResourceManager.")
			res := c.protoResourceman.HandleResourceManagerResponse(wrapper)

			switch res.(type) {

			case *proto3Resourceman.JobResponse:
				response := res.(*proto3Resourceman.JobResponse)

				if response.Err != nil {
					c.logger.Sugar().Error("Error response from the resourcemanager: ", zap.Error(response.Err))
					return
				}
			case *proto3Resourceman.ProgressUpdate:

				response := res.(*proto3Resourceman.ProgressUpdate)

				if response == nil {
					c.logger.Info("Received nil progress update from the resource manager.")
					c.logger.Info("Closing connection with the resource manager.")
					c.protoResourceman.MsgHandler().Close()
					return

				}

				updateType := response.ProgressUpdateType

				c.logger.Info("Received progress update from the resource manager.")

				c.logger.Info("Progress update type: ", zap.String("updateType", updateType))
				if response.TaskProgress != nil {
					c.logger.Info("Progress update count: ", zap.String("Phase", response.TaskProgress.TaskType))
					c.logger.Info("Progress update count: ", zap.Int32("count", response.TaskProgress.TaskNumber))
					c.logger.Info("Progress update count: ", zap.Int32("Total", response.TaskProgress.TotalTasks))
				}

				fmt.Println()
				fmt.Println()
				fmt.Println()

			}

		case nil:
			c.logger.Info("Received nil message from the ResourceManager.")
			return
		}

	}

}

// DialResourceManager dials the resource manager and sets up the connection with it.
func (c *Client) DialResourceManager() (err error) {
	server := c.serverPort

	conn, err := net.Dial("tcp", server)

	if err != nil {
		c.logger.Fatal("There was an error connecting to the host.")
	}
	msgHandler := messagesResourceman.NewMessageHandler(conn)
	proto := proto3Resourceman.NewProtoHandler(msgHandler, c.logger)
	c.SetProtoResourceman(proto)
	c.Conn(conn)

	return
}

// HandleMRJob handles the map reduce job.
func (c *Client) HandleMRJob(inputFile, outputFile []string, reducerCount []int32, jobBinaryPath [][]byte) {
	c.logger.Info("Handling MR job.")
	c.protoResourceman.HandleMRJobRequest(inputFile, outputFile, reducerCount, jobBinaryPath)
}
