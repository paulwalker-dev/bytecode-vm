package main

import (
	"flag"
	"fmt"
	"github.com/paulwalker-dev/os/host/pkg/emulator"
	"log"
	"os"
)

func main() {
	emu := emulator.New()
	emu.Reset()

	flag.Parse()
	mem, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	for i, b := range mem {
		emu.SetByte(uint16(i), b)
	}

	emu.Start()
	emu.Wait()

	for i := 1; i <= emulator.RegisterCount; i++ {
		fmt.Printf("%#.2x: %#.04x\n", i, emu.GetRegister(byte(i)))
	}
}
