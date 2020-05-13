build: clean go-build package

build-arm-linux: clean go-build-arm-linux package

build-arm64-linux: clean go-build-arm64-linux package

.PHONY: go-build
go-build:
	go build

.PHONY: go-build-arm-linux
go-build-arm-linux:
	# Needs ARM compiler (gcc-arm-linux-gnueabihf)
	CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build

.PHONY: go-build-arm64-linux
go-build-arm64-linux:
	# Needs ARM64 compiler (gcc-aarch64-linux-gnu)
	CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build

.PHONY: package
package:
	zip work-log.zip \
		work-log \
		config/config.ini.example \
		resources/css/*.* \
		resources/font/*.* \
		resources/img/*.* \
		resources/js/*.* \
		scripts/db/*.* \
		templates/*.*

.PHONY: clean
clean:
	rm -f work-log
	rm -f work-log.zip