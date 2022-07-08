package main

import local "github.com/lamhai1401/rpc-ex/pubsub"

func main() {
	// go runServer()
	// go runGRPCServer()
	go local.Run()

	go local.RunClientSub()
	local.RunClientPub()

	select {}
}
