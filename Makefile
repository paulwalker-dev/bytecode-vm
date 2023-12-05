all:
	@tup

run: all
	./build/host/bin/emu ./build/boot/boot.bin

.PHONY: all
