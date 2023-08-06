package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
)

type NodeManager struct {
	NodeManagerPort     string
	ResourceManagerPort string
	dir                 string
	logger              *zap.Logger
}

func NewNodeManager(NodeManagerPort, ResourceManagerPort string, logger *zap.Logger, dir string) *NodeManager {
	newNodeManager := &NodeManager{
		NodeManagerPort:     NodeManagerPort,
		ResourceManagerPort: ResourceManagerPort,
		logger:              logger,
		dir:                 dir,
	}
	return newNodeManager
}

func (nm *NodeManager) ListenOnPort() (nodeManListener net.Listener, resourceManListener net.Listener) {
	nm.logger.Info("Node Manager Started")

	var err error
	nodeManListener, err = net.Listen("tcp", ":"+nm.NodeManagerPort)

	if err != nil {
		nm.logger.Fatal("Error listening on NodeManagerPort " + nm.NodeManagerPort)
	}

	resourceManListener, err = net.Listen("tcp", ":"+nm.ResourceManagerPort)
	if err != nil {
		nm.logger.Fatal("Error listening on ResourceManagerPort " + nm.ResourceManagerPort)
	}

	return nodeManListener, resourceManListener

}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run node_manager.go <Directory that hosts the node-manager>")
		return
	}
	logger := Logger()

	//TODO: Get NodeManagerPort from command line arguments or config file
	nm := NewNodeManager("23098", "23099", logger, os.Args[1])
	nodeManListener, resourceManListener := nm.ListenOnPort()
	go nm.AcceptResourceManagerConnections(resourceManListener)
	go nm.AcceptNodeManagerConnections(nodeManListener)

	select {}
}

func Logger() *zap.Logger {

	file, err := os.OpenFile("logfileN.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Create a logger that writes to the file
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter := zapcore.AddSync(file)
	fileLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, fileLevel)

	// Create a logger that writes to the console
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriter := zapcore.Lock(os.Stdout)
	consoleLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, consoleLevel)

	// Create a final logger that writes to both the file and console
	logger := zap.New(zapcore.NewTee(fileCore, consoleCore))
	return logger
}
