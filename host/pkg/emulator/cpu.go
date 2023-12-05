package emulator

import (
	"fmt"
	"sync"
)

const RegisterCount = 16

type CPU interface {
	Running() bool
	Wait()

	Start()
	Pause()
	Reset()

	GetRegister(id int) uint16
	SetRegister(id int, value uint16)
}

type cpu struct {
	m   sync.RWMutex
	ram RAM

	running bool
	signal  chan bool

	registers [RegisterCount]uint16
	pc, sp    uint16
}

const (
	_ int = -iota
	PC
	SP
	lowest
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
	NOP byte = iota
	ADD
	IMM
	LOAD
	STORE
	MOV
	JMP
	BEZ
	SUB
	HALT byte = 0xFF
)

func (p *cpu) loop() {
	for p.Running() {
		pc := p.GetRegister(PC)
		opcode := p.ram.GetByte(pc)
		start := pc
		switch opcode {
		case NOP:
		case ADD:
			dest := int(p.ram.GetByte(pc + 1))
			src1 := p.GetRegister(int(p.ram.GetByte(pc + 2)))
			src2 := p.GetRegister(int(p.ram.GetByte(pc + 3)))
			p.SetRegister(dest, src1+src2)
			pc += 3
		case IMM:
			dest := int(p.ram.GetByte(pc + 1))
			low := uint16(p.ram.GetByte(pc + 2))
			high := uint16(p.ram.GetByte(pc + 3))
			p.SetRegister(dest, high<<8^low)
			pc += 3
		case LOAD:
			dest := int(p.ram.GetByte(pc + 1))
			addr := p.GetRegister(int(p.ram.GetByte(pc + 2)))
			low := uint16(p.ram.GetByte(addr))
			high := uint16(p.ram.GetByte(addr + 1))
			p.SetRegister(dest, high<<8|low)
			pc += 2
		case STORE:
			addr := p.GetRegister(int(p.ram.GetByte(pc + 1)))
			data := p.GetRegister(int(p.ram.GetByte(pc + 2)))
			p.ram.SetByte(addr, uint8(data))
			p.ram.SetByte(addr+1, uint8(data>>8))
			pc += 2
		case MOV:
			dest := int(p.ram.GetByte(pc + 1))
			src := int(p.ram.GetByte(pc + 2))
			p.SetRegister(dest, p.GetRegister(src))
			pc += 2
		case JMP:
			pc += 1
		case BEZ:
			data := p.GetRegister(int(p.ram.GetByte(pc + 1)))
			if data == 0 {
				opcode = JMP
			}
			pc += 2
		case SUB:
			dest := int(p.ram.GetByte(pc + 1))
			src1 := p.GetRegister(int(p.ram.GetByte(pc + 2)))
			src2 := p.GetRegister(int(p.ram.GetByte(pc + 3)))
			p.SetRegister(dest, src1-src2)
			pc += 3
		case HALT:
			p.Pause()
		}
		pc++

		fmt.Printf("%#.4x:", start)
		for ; start < pc; start++ {
			fmt.Printf(" %.2x", p.ram.GetByte(start))
		}
		fmt.Println()

		if opcode == JMP {
			pc = p.GetRegister(int(p.ram.GetByte(pc - 1)))
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
	for i := lowest; i < RegisterCount; i++ {
		p.SetRegister(i+1, 0)
	}
}

func (p *cpu) getRegisterPointer(id int) *uint16 {
	if id <= 0 {
		switch id {
		case PC:
			return &p.pc
		case SP:
			return &p.sp
		}
	} else if id <= RegisterCount {
		return &p.registers[id-1]
	}
	buf := uint16(0)
	return &buf
}

func (p *cpu) GetRegister(id int) uint16 {
	p.m.RLock()
	defer p.m.RUnlock()

	return *p.getRegisterPointer(id)
}

func (p *cpu) SetRegister(id int, value uint16) {
	p.m.Lock()
	defer p.m.Unlock()

	*p.getRegisterPointer(id) = value
}
