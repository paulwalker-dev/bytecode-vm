all:
	@tup

run: all
	./host/bin/emu/emu -disk ./boot/mbr.bin

.PHONY: all
