package devices

import (
	"os"

	"github.com/paulwalker-dev/bytecode-vm/host/pkg/emulator"
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
	d.block = uint16(value)<<8 | d.block>>8
}

func (d *disk) Reset() {}
