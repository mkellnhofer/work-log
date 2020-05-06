rm work-log
rm work-log-linux.zip

go build

zip work-log-linux.zip work-log
zip -u work-log-linux.zip config/config.ini.example
zip -u work-log-linux.zip resources/css/*.* resources/font/*.* resources/img/*.* resources/js/*.*
zip -u work-log-linux.zip scripts/db/*.*
zip -u work-log-linux.zip templates/*.*