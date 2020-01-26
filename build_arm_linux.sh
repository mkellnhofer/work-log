rm work-log
rm work-log-arm-linux.zip

# Needs ARM compiler (gcc-arm-linux-gnueabihf)
CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build

zip work-log-arm-linux.zip work-log
zip -u work-log-arm-linux.zip config/config.ini.example
zip -u work-log-arm-linux.zip resources/css/*.* resources/img/*.* resources/js/*.*
zip -u work-log-arm-linux.zip scripts/db/*.*
zip -u work-log-arm-linux.zip templates/*.*