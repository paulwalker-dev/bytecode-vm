all:
	@mkdir -p build
	@touch build/tup.config
	@tup

run: all
	./build/host/bin/emu -disk ./build/sys/sys.img

.PHONY: clean
clean:
	rm -rf build
	git clean -fXd
