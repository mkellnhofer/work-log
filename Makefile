sass:
	sass \
		scss/bootstrap-custom.scss \
		resources/static/web/css/bootstrap.min.css \
		--style=compressed \
		--no-source-map \
		--silence-deprecation=mixed-decls

VERSION ?=
PLATFORM ?=

PACKAGE_NAME = work-log$(if $(VERSION),-$(VERSION),)$(if $(PLATFORM),-$(PLATFORM),).zip

build: clean swagger-spec go-build package

build-amd64-linux: PLATFORM = amd64-linux
build-amd64-linux: clean swagger-spec go-build-amd64-linux package

build-arm-linux: PLATFORM = arm-linux
build-arm-linux: clean swagger-spec go-build-arm-linux package

build-arm64-linux: PLATFORM = arm64-linux
build-arm64-linux: clean swagger-spec go-build-arm64-linux package

.PHONY: clean
clean:
	rm -f work-log
	rm -f $(PACKAGE_NAME)

.PHONY: swagger-spec
swagger-spec:
	# Needs go-swagger
	SWAGGER_GENERATE_EXTENSION="false" go-swagger generate spec \
		-o resources/static/swagger-ui/swagger.json

.PHONY: go-build
go-build:
	go build -o work-log cmd/init.go cmd/main.go

.PHONY: go-build-amd64-linux
go-build-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o work-log cmd/init.go cmd/main.go

.PHONY: go-build-arm-linux
go-build-arm-linux:
	GOOS=linux GOARCH=arm GOARM=6 go build -o work-log cmd/init.go cmd/main.go

.PHONY: go-build-arm64-linux
go-build-arm64-linux:
	GOOS=linux GOARCH=arm64 go build -o work-log cmd/init.go cmd/main.go

.PHONY: package
package:
	zip $(PACKAGE_NAME) \
		work-log \
		config/config.ini.example \
		resources/db/*.* \
		resources/localizations/*.* \
		resources/static/swagger-ui/*.* \
		resources/static/web/css/*.* \
		resources/static/web/font/*.* \
		resources/static/web/img/*.* \
		resources/static/web/js/*.*
