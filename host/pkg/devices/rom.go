package devices

import "github.com/paulwalker-dev/os/host/pkg/emulator"

type ROM interface {
	emulator.Device
	SetContents(value []byte)
}

type rom struct {
	contents []byte
}

func NewROM(contents []byte) ROM {
	return &rom{contents: contents}
}

func (d *rom) GetByte(addr uint8) byte {
	if int(addr) < len(d.contents) {
		return d.contents[addr]
	}
	return 0
}

func (d *rom) SetByte(addr uint8, value byte) {}

func (d *rom) Reset() {}

func (d *rom) SetContents(value []byte) {
	d.contents = value
}
