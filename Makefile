build: clean swagger-spec go-build package

build-arm-linux: clean swagger-spec go-build-arm-linux package

build-arm64-linux: clean swagger-spec go-build-arm64-linux package

.PHONY: clean
clean:
	rm -f work-log
	rm -f work-log.zip

.PHONY: swagger-spec
swagger-spec:
	# Needs go-swagger
	SWAGGER_GENERATE_EXTENSION="false" go-swagger generate spec -o static/swagger-ui/swagger.json

.PHONY: go-build
go-build:
	go build -o work-log cmd/init.go cmd/main.go

.PHONY: go-build-arm-linux
go-build-arm-linux:
	GOOS=linux GOARCH=arm GOARM=6 go build -o work-log cmd/init.go cmd/main.go

.PHONY: go-build-arm64-linux
go-build-arm64-linux:
	GOOS=linux GOARCH=arm64 go build -o work-log cmd/init.go cmd/main.go

.PHONY: package
package:
	zip work-log.zip \
		work-log \
		config/config.ini.example \
		config/localizations/*.* \
		scripts/db/*.* \
		static/resources/css/*.* \
		static/resources/font/*.* \
		static/resources/img/*.* \
		static/resources/js/*.* \
		static/swagger-ui/*.*