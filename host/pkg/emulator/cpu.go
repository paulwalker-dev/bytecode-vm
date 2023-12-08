package emulator

import "sync"

const RegisterCount = 16

type CPU interface {
	Running() bool
	Wait()

	Start()
	Pause()
	Reset()

	GetRegister(id byte) uint16
	SetRegister(id byte, value uint16)
}

type cpu struct {
	m   sync.RWMutex
	ram RAM

	running bool
	signal  chan bool

	registers [RegisterCount]uint16
}

const (
	PC byte = RegisterCount - iota
	SP
	BUF
)

func newCPU(ram RAM) CPU {
	return &cpu{
		ram:    ram,
		signal: make(chan bool),
	}
}

func (p *cpu) Running() bool {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.running
}

func (p *cpu) Wait() {
	if p.Running() {
		<-p.signal
	}
}

func (p *cpu) Start() {
	if p.Running() {
		return
	}

	p.m.Lock()
	p.running = true
	p.m.Unlock()
	go p.loop()
}

const (
	NOP   byte = 0x00
	IMM   byte = 0x01
	MOV   byte = 0x02
	STORE byte = 0x10
	LOAD  byte = 0x11
	ADD   byte = 0x20
	SUB   byte = 0x21
	XOR   byte = 0x22
	NOT   byte = 0x23
	JMP   byte = 0x30
	BEZ   byte = 0x31
	HALT  byte = 0xFF
)

func (p *cpu) loop() {
	for p.Running() {
		pc := p.GetRegister(PC)
		opcode := p.ram.GetByte(pc)
		switch opcode {
		case NOP:
		case IMM:
			dest := p.ram.GetByte(pc + 1)
			low := uint16(p.ram.GetByte(pc + 2))
			high := uint16(p.ram.GetByte(pc + 3))
			p.SetRegister(dest, high<<8^low)
			pc += 3
		case MOV:
			dest := p.ram.GetByte(pc + 1)
			src := p.ram.GetByte(pc + 2)
			p.SetRegister(dest, p.GetRegister(src))
			pc += 2
		case STORE:
			addr := p.GetRegister(p.ram.GetByte(pc + 1))
			data := p.GetRegister(p.ram.GetByte(pc + 2))
			p.ram.SetByte(addr, uint8(data))
			p.ram.SetByte(addr+1, uint8(data>>8))
			pc += 2
		case LOAD:
			dest := p.ram.GetByte(pc + 1)
			addr := p.GetRegister(p.ram.GetByte(pc + 2))
			low := uint16(p.ram.GetByte(addr))
			high := uint16(p.ram.GetByte(addr + 1))
			p.SetRegister(dest, high<<8|low)
			pc += 2
		case ADD:
			dest := p.ram.GetByte(pc + 1)
			src1 := p.GetRegister(p.ram.GetByte(pc + 2))
			src2 := p.GetRegister(p.ram.GetByte(pc + 3))
			p.SetRegister(dest, src1+src2)
			pc += 3
		case SUB:
			dest := p.ram.GetByte(pc + 1)
			src1 := p.GetRegister(p.ram.GetByte(pc + 2))
			src2 := p.GetRegister(p.ram.GetByte(pc + 3))
			p.SetRegister(dest, src1-src2)
			pc += 3
		case XOR:
			dest := p.ram.GetByte(pc + 1)
			src1 := p.GetRegister(p.ram.GetByte(pc + 2))
			src2 := p.GetRegister(p.ram.GetByte(pc + 3))
			p.SetRegister(dest, src1^src2)
			pc += 3
		case NOT:
			dest := p.ram.GetByte(pc + 1)
			src := p.GetRegister(p.ram.GetByte(pc + 2))
			p.SetRegister(dest, 0xFFFF^src)
			pc += 2
		case JMP:
			pc += 1
		case BEZ:
			data := p.GetRegister(p.ram.GetByte(pc + 1))
			if data == 0 {
				opcode = JMP
			}
			pc += 2
		case HALT:
			p.Pause()
		}
		pc++

		if opcode == JMP {
			pc = p.GetRegister(p.ram.GetByte(pc - 1))
		}

		p.SetRegister(PC, pc)
	}
	p.signal <- true
}

func (p *cpu) Pause() {
	p.m.Lock()
	defer p.m.Unlock()

	p.running = false
}

func (p *cpu) Reset() {
	for i := 0; i < RegisterCount; i++ {
		p.SetRegister(byte(i+1), 0)
	}
}

func (p *cpu) getRegisterPointer(id byte) *uint16 {
	if 0 < id && id <= RegisterCount {
		return &p.registers[id-1]
	}
	buf := uint16(0)
	return &buf
}

func (p *cpu) GetRegister(id byte) uint16 {
	p.m.RLock()
	defer p.m.RUnlock()

	return *p.getRegisterPointer(id)
}

func (p *cpu) SetRegister(id byte, value uint16) {
	p.m.Lock()
	defer p.m.Unlock()

	*p.getRegisterPointer(id) = value
}
