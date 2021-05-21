package main

func main() {
	var b bus

	r := createRAM()
	b.connect(0x0000, 0xFFFF, r)

	createCPU(&b)
}
