package main

func main() {
	defer cleanup()
	in, out := choosePorts()
	_ = relay(in, out)
	select {}
}
