package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/paulwalker-dev/os/host/pkg/devices"
	"github.com/paulwalker-dev/os/host/pkg/emulator"
)

//go:embed bios.bin
var bios []byte

func main() {
	disk := flag.String("disk", "", "The disk to use")
	flag.Parse()

	emu := emulator.New()
	emu.Reset()
	emu.AssignDevice(0x00, devices.NewROM(bios))

	if *disk != "" {
		emu.AssignDevice(0x10, devices.NewDisk(*disk))
	}

	emu.Start()
	emu.Wait()

	for i := 1; i <= emulator.RegisterCount; i++ {
		fmt.Printf("%#.2x: %#.04x\n", i, emu.GetRegister(byte(i)))
	}
	fmt.Printf("RET: %#.4x\n", emu.GetByte(emu.GetRegister(emulator.SP)-2))
}
