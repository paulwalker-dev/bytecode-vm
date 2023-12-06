include "include/cpu.inc"
org 0x0000

	jmpl entry
	nop
	halt

; 2 Disk device
; 3 Block id
; 4 Destination
loadblock:
	imm 5, 0xfe
	store 2, 3
.loop:
	load 1, 2
	store 4, 1
	imm 1, 2
	add 2, 2, 1
	add 4, 4, 1
	sub 5, 5, 1
	imm 1, .done
	bez 5, 1
	jmpl .loop
.done:
	ret

entry:
	imm sp, 0xff00
	imm 2, 0x1000
	imm 3, 0x0000
	imm 4, 0x0100
	call loadblock
	imm 2, 0x1000
	imm 3, 0x0001
	call loadblock
	jmpl 0x0100
	halt

repeat 0x100-($-$$)
	nop
end repeat
