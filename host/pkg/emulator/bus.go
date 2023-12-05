package emulator

import "sync"

type Device interface {
	GetByte(addr uint8) byte
	SetByte(addr uint8, value byte)

	Reset()
}

type RAM interface {
	GetByte(addr uint16) byte
	SetByte(addr uint16, value byte)
}

type Bus interface {
	GetByte(addr uint16) byte
	SetByte(addr uint16, value byte)

	AssignDevice(block uint8, device Device)
	GetDevice(block uint8) Device
	RemoveDevice(block uint8)

	Reset()
}

type bus struct {
	m sync.RWMutex

	ram     [0x10000]byte
	devices map[uint8]Device
}

func newBus() Bus {
	return &bus{}
}

func (b *bus) GetByte(addr uint16) byte {
	b.m.RLock()
	defer b.m.RUnlock()

	block := uint8(addr / 0x100)
	if device, present := b.devices[block]; present {
		return device.GetByte(uint8(addr % 0x100))
	}
	return b.ram[addr]
}

func (b *bus) SetByte(addr uint16, value byte) {
	b.m.Lock()
	defer b.m.Unlock()

	block := uint8(addr / 0x100)
	if device, present := b.devices[block]; present {
		device.SetByte(uint8(addr%0x100), value)
	}
	b.ram[addr] = value
}

func (b *bus) AssignDevice(block uint8, device Device) {
	b.m.Lock()
	defer b.m.Unlock()

	if device != nil {
		b.devices[block] = device
	}
}

func (b *bus) GetDevice(block uint8) Device {
	b.m.RLock()
	defer b.m.RUnlock()

	return b.devices[block]
}

func (b *bus) RemoveDevice(block uint8) {
	b.m.Lock()
	defer b.m.Unlock()

	if _, present := b.devices[block]; present {
		delete(b.devices, block)
	}
}

func (b *bus) Reset() {
	b.m.Lock()
	defer b.m.Unlock()

	for _, device := range b.devices {
		device.Reset()
	}
	for i := 0; i < 0x10000; i++ {
		b.ram[i] = 0
	}
}
