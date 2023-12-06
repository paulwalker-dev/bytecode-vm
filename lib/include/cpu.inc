PC? equ 16
SP? equ 15
BUF? equ 14

macro NOP?
	db 0x00
end macro

macro ADD? dest,src1,src2
	db 0x01, dest, src1, src2
end macro

macro IMM? dest,val
	db 0x02, dest
	emit 2: val
end macro

macro LOAD? dest,addr
	db 0x03, dest, addr
end macro

macro STORE? addr,src
	db 0x04, addr, src
end macro

macro MOV? dest,src
	db 0x05, dest, src
end macro

macro JMP? dest
	db 0x06, dest
end macro

macro BEZ? dest,src
	db 0x07, dest, src
end macro

macro SUB? dest,src1,src2
	db 0x08, dest, src1, src2
end macro

macro HALT?
	db 0xFF
end macro

macro PUSH? src
	STORE SP, src
	IMM BUF, 2
	ADD SP, SP, BUF
end macro

macro POP? dest
	IMM BUF, 2
	SUB SP, SP, BUF
	LOAD dest, SP
end macro

macro CALL? addr
	local callback
	IMM BUF, callback
	PUSH BUF
	IMM BUF, addr
	JMP BUF
	callback:
end macro

macro RET?
	POP BUF
	JMP BUF
end macro
