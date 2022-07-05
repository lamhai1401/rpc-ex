package main

func main() {
	// go runServer()
	go runServeCallLient()

	// runClient()
	runRPCClient()

	select {}
}
