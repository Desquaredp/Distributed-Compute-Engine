package main

import (
	"net"
	"sort"
	messagesStorage "src/messages/client_storage"
	proto3Storage "src/proto/client_storage"
	proto3 "src/proto/controller_client"
	"strconv"
	"strings"
)

//var sem = make(chan struct{}, 1) // allow up to 1 goroutines at once

func (c *Client) DispatchFile(res proto3.ResponseInterface) {

	fragments := res.(*proto3.PlanResponse).FragmentLayout
	fragments = c.SortFragments(fragments)
	c.file.SetFragmentLayout(fragments)

	//sort the fragments by fragId

	for _, frag := range fragments {
		//TODO: can be done in parallel to the main thread
		c.logger.Info("Dispatching fragment: " + frag.FragmentId)
		c.DispatchFragment(frag)
	}

	c.logger.Info("All fragments dispatched")

}

func (c *Client) SortFragments(fragments []proto3.FragmentInfo) (sortedFragments []proto3.FragmentInfo) {
	// Sort the fragments based on their fragment number
	sort.Slice(fragments, func(i, j int) bool {
		f1 := fragments[i]
		f2 := fragments[j]

		// Extract the fragment numbers from the fragment names
		f1Num, _ := strconv.Atoi(strings.Split(f1.FragmentId, "_")[1])
		f2Num, _ := strconv.Atoi(strings.Split(f2.FragmentId, "_")[1])

		return f1Num < f2Num
	})

	return fragments
}

func (c *Client) DispatchFragment(frag proto3.FragmentInfo) {
	//sem <- struct{}{}        // acquire semaphore
	//defer func() { <-sem }() // release semaphore

	nodes := frag.StorageNodes

	//TODO: if no nodes are available, then wait for a node to be available or panic
	for _, node := range nodes {
		err := c.DispatchToNode(frag, node)
		if err != nil {
			c.logger.Sugar().Error("There was an error dispatching to node: ", node.NodeId)
			continue
		} else {
			break
		}
	}

}

func (c *Client) DispatchToNode(frag proto3.FragmentInfo, node proto3.StorageNodeInfo) (err error) {

	conn, err := net.Dial("tcp", node.Host+":"+node.Port)
	if err != nil {
		c.logger.Error("There was an error connecting to the host.")
		return
	}

	msgHandler := messagesStorage.NewMessageHandler(conn)
	proto := proto3Storage.NewProtoHandler(msgHandler, c.logger, c.file.Dir())
	proto.SetFileHandler(c.file)
	//might break if too many goroutines are created

	err = proto.HandleFilePutRequest(frag.FragmentId)
	if err != nil {
		c.logger.Error("There was an error handling the file put request.")
		return
	}

	c.handleStorageResponse(proto)
	return
}

func (c *Client) handleStorageResponse(proto *proto3Storage.ProtoHandler) {

	defer proto.MsgHandler().Close()

	for {

		wrapper, err := proto.MsgHandler().ServerResponseReceive()
		if err != nil {
			c.logger.Error("There was an error receiving the server response.")
			return
		}

		switch wrapper.Response.(type) {

		default:
			proto.HandleResponse(wrapper)

		case nil:
			return
		}
	}
}
