rm work-log
rm work-log-arm64-linux.zip

# Needs ARM64 compiler (gcc-aarch64-linux-gnu)
CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build

zip work-log-arm64-linux.zip work-log
zip -u work-log-arm64-linux.zip config/config.ini.example
zip -u work-log-arm64-linux.zip resources/css/*.* resources/img/*.* resources/js/*.*
zip -u work-log-arm64-linux.zip scripts/db/*.*
zip -u work-log-arm64-linux.zip templates/*.*