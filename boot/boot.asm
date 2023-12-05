	include "include/cpu.inc"

	imm 16, 0xFF00

start:
	imm 2, 2
	imm 3, 8
	call power
	halt

mult:
	imm 1, 0
	push 3
.loop:
	add 1, 1, 2
	imm 15, 1
	sub 3, 3, 15
	imm 15, .done
	bez 3, 15
	imm 15, .loop
	jmp 15
.done:
	pop 3
	ret

power:
	imm 15, .zero
	bez 3, 15
	push 4
	push 3
	mov 4, 3
	mov 3, 2
	imm 2, 1
.loop:
	imm 15, .done
	bez 4, 15
	call mult
	mov 2, 1
	imm 15, 1
	sub 4, 4, 15
	imm 15, .loop
	jmp 15
.zero:
	imm 1, 1
	ret
.done:
	mov 2, 3
	pop 3
	pop 4
	ret
