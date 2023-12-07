package devices

import (
	"github.com/paulwalker-dev/os/host/pkg/emulator"
	"os"
)

type Disk interface {
	emulator.Device
	Open(path string)
}

type disk struct {
	filePath string
	block    uint16
}

func NewDisk(path string) Disk {
	return &disk{filePath: path}
}

func (d *disk) Open(file string) {
	d.filePath = file
}

func (d *disk) GetByte(addr uint8) byte {
	fi, err := os.OpenFile(d.filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0xfe
	}
	defer fi.Close()
	ret := make([]byte, 1)
	_, err = fi.ReadAt(ret, int64((d.block*0x100)+uint16(addr)))
	if err != nil {
		return 0xff
	}
	return ret[0]
}

func (d *disk) SetByte(addr uint8, value byte) {
	if addr == 0 {
		d.block = uint16(value)
	}
}

func (d *disk) Reset() {}
