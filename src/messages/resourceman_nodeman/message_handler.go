package resourceman_nodeman

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

type MessageHandler struct {
	conn net.Conn
}

func NewMessageHandler(conn net.Conn) *MessageHandler {
	m := &MessageHandler{
		conn: conn,
	}

	return m
}

func (m *MessageHandler) readN(buf []byte) error {
	bytesRead := uint64(0)
	for bytesRead < uint64(len(buf)) {
		n, err := m.conn.Read(buf[bytesRead:])
		if err != nil {
			return err
		}
		bytesRead += uint64(n)
	}
	return nil
}

func (m *MessageHandler) writeN(buf []byte) error {
	bytesWritten := uint64(0)
	for bytesWritten < uint64(len(buf)) {
		n, err := m.conn.Write(buf[bytesWritten:])
		if err != nil {
			return err
		}
		bytesWritten += uint64(n)
	}
	return nil
}

func (m *MessageHandler) ResourceManagerRequestSend(wrapper *ResourceManagerTaskRequest) error {
	serialized, err := proto.Marshal(wrapper)
	if err != nil {
		return err
	}

	prefix := make([]byte, 8)
	binary.LittleEndian.PutUint64(prefix, uint64(len(serialized)))
	m.writeN(prefix)
	m.writeN(serialized)

	return nil
}

func (m *MessageHandler) NodeManagerResponseSend(wrapper *NodeManagerTaskResponse) error {
	serialized, err := proto.Marshal(wrapper)
	if err != nil {
		return err
	}

	prefix := make([]byte, 8)
	binary.LittleEndian.PutUint64(prefix, uint64(len(serialized)))
	m.writeN(prefix)
	m.writeN(serialized)

	return nil

}

func (m *MessageHandler) ResourceManagerRequestReceive() (*ResourceManagerTaskRequest, error) {
	prefix := make([]byte, 8)
	m.readN(prefix)

	payloadSize := binary.LittleEndian.Uint64(prefix)
	payload := make([]byte, payloadSize)
	m.readN(payload)

	wrapper := &ResourceManagerTaskRequest{}
	err := proto.Unmarshal(payload, wrapper)
	return wrapper, err
}

func (m *MessageHandler) NodeManagerResponseReceive() (*NodeManagerTaskResponse, error) {
	prefix := make([]byte, 8)
	m.readN(prefix)

	payloadSize := binary.LittleEndian.Uint64(prefix)
	payload := make([]byte, payloadSize)
	m.readN(payload)

	wrapper := &NodeManagerTaskResponse{}
	err := proto.Unmarshal(payload, wrapper)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	return wrapper, err
}

func (m *MessageHandler) Close() {
	m.conn.Close()
}
