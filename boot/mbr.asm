	include "include/cpu.inc"

org 0x0100

	jmpl start
	nop
	halt

;   1 Buffer
;   2 Disk device
;   3 Block id
; * 4 Destination
;   5 Counter
loadblock:
	push 1
	push 2
	push 5
	call 0x0008
	pop 5
	pop 2
	pop 1
	ret

blocks = (eof-0x0200)/0x100
start:
	imm sp, 0xff00
	imm 2, 0x1000
	imm 3, 0x0002
	imm 4, 0x0300
	imm 5, blocks
.loop:
	call loadblock
	imm 1, 1
	add 3, 3, 1
	sub 5, 5, 1
	imm 1, stage2
	bez 5, 1
	jmpl .loop

repeat 0x1fe-($-$$)
	nop
end repeat

	db 0x55, 0xaa

stage2:
	imm 1, 0xffff
	halt
eof:
