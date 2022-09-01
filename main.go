package main

func main() {
	defer cleanup()
	in, out := choosePorts()
	scale := newScale()
	scale.addTamperament("*", "B", -50)
	scale.addTamperament("5", "Gb", -50)
	_ = relay(in, out, scale)
	select {}
}
