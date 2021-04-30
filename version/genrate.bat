cd "%~dp0"
call make_version.bat ./  ./version.h

cd "%~dp0"
windres.exe -i yc.rc -o yc.syso
 
rem cd ..
rem call go build -ldflags "-H windowsgui"

pause