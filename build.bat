set dir=%~dp0
cd /D %dir%

cd version

del /F /S /Q  version.h

call genrate.bat

copy   yc.syso   ..\

cd ../

go build

pause