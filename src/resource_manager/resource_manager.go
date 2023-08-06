package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	messagesController "src/messages/resourceman_controller"
	"strconv"
)

// ResourceManager is the main struct for the resource manager
type ResourceManager struct {
	logger         *zap.Logger
	controllerAddr string
}

func NewResourceManager(logger *zap.Logger) *ResourceManager {
	newResourceManager := &ResourceManager{
		logger: logger,
	}
	return newResourceManager

}

func main() {

	logger := Logger()

	rm := NewResourceManager(logger)

	if len(os.Args) < 4 {
		rm.logger.Warn("Command line arguments are not correct")
		rm.logger.Fatal("Usage: go run resource_manager.go <Controller host:port> <Node Manager Facing Port> <Client facing Port>")
		os.Exit(1)
	}

	rm.logger.Info("Resource Manager Started")

	controllerAddr := os.Args[1]
	rm.controllerAddr = controllerAddr
	//rm.DialController(controllerAddr)

	nodeManagerPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		rm.logger.Fatal("Node Manager Port is not an integer")
		os.Exit(1)
	}
	clientFacingPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		rm.logger.Fatal("Client Facing Port is not an integer")
		os.Exit(1)
	}

	fmt.Println("Controller Address: ", controllerAddr)
	fmt.Println("Node Manager Port: ", nodeManagerPort)
	fmt.Println("Client Facing Port: ", clientFacingPort)

	_, err = rm.ListenOnPort(nodeManagerPort)
	if err != nil {
		rm.logger.Fatal(err.Error())
	}

	clientFacingListener, err := rm.ListenOnPort(clientFacingPort)
	if err != nil {
		rm.logger.Fatal(err.Error())
	}

	go rm.acceptClientConnections(clientFacingListener)

	select {}

}

// ListenOnPort listens on a specified port
func (rm *ResourceManager) ListenOnPort(port int) (net.Listener, error) {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		rm.logger.Error(err.Error())
		return nil, err
	}
	return listener, nil
}

func (rm *ResourceManager) DialController() (msgHandler *messagesController.MessageHandler) {

	conn, err := net.Dial("tcp", rm.controllerAddr)
	if err != nil {
		rm.logger.Sugar().Fatal("There was error while dialing to the controller: \n", err.Error())
	}

	msgHandler = messagesController.NewMessageHandler(conn)
	//proto := proto3Controller.NewProtoHandler(msgHandler, rm.logger)
	//rm.protoController = proto
	return
}

func Logger() *zap.Logger {

	file, err := os.OpenFile("logfileR.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
