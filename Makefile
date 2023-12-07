all:
	@tup

run: all
	./host/bin/emu/emu -disk ./sys/sys.img

.PHONY: all
