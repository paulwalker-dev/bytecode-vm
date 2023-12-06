org 0x0100

	include "include/cpu.inc"

start:
	imm 1, 0xffff
	halt

loadblock:
	push 1
	push 2
	push 4
	push 5
	call 0x000e
	pop 5
	pop 4
	pop 2
	pop 1
	ret

repeat 0x200-($-$$)
	nop
end repeat
