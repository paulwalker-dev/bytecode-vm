	include "include/cpu.inc"

	imm sp, 0xFF00
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
	imm buf, 1
	sub 3, 3, buf
	imm buf, .done
	bez 3, buf
	imm buf, .loop
	jmp buf
.done:
	pop 3
	ret

power:
	imm buf, .zero
	bez 3, buf
	push 4
	push 3
	mov 4, 3
	mov 3, 2
	imm 2, 1
.loop:
	imm buf, .done
	bez 4, buf
	call mult
	mov 2, 1
	imm buf, 1
	sub 4, 4, buf
	imm buf, .loop
	jmp buf
.zero:
	imm 1, 1
	ret
.done:
	mov 2, 3
	pop 3
	pop 4
	ret
