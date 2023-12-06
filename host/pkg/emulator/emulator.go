package emulator

type Emulator interface {
	CPU
	Bus
}

type emulator struct {
	cpu CPU
	bus Bus
}

func New() Emulator {
	emu := &emulator{}

	emu.bus = newBus()
	emu.cpu = newCPU(emu.bus)

	emu.Reset()
	return emu
}

func (emu *emulator) Running() bool {
	return emu.cpu.Running()
}

func (emu *emulator) Wait() {
	emu.cpu.Wait()
}

func (emu *emulator) Start() {
	emu.cpu.Start()
}

func (emu *emulator) Pause() {
	emu.cpu.Pause()
}

func (emu *emulator) Reset() {
	emu.cpu.Pause()
	emu.bus.Reset()
	emu.cpu.Reset()
}

func (emu *emulator) GetRegister(id byte) uint16 {
	return emu.cpu.GetRegister(id)
}

func (emu *emulator) SetRegister(id byte, value uint16) {
	emu.cpu.SetRegister(id, value)
}

func (emu *emulator) GetByte(addr uint16) byte {
	return emu.bus.GetByte(addr)
}

func (emu *emulator) SetByte(addr uint16, value byte) {
	emu.bus.SetByte(addr, value)
}

func (emu *emulator) AssignDevice(block uint8, device Device) {
	emu.bus.AssignDevice(block, device)
}

func (emu *emulator) GetDevice(block uint8) Device {
	return emu.bus.GetDevice(block)
}

func (emu *emulator) RemoveDevice(block uint8) {
	emu.bus.RemoveDevice(block)
}
