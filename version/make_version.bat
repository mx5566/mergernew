@ECHO OFF
cd /d %1
if exist %2 del /q %2
for /f "delims=" %%i in ('git rev-list --count HEAD') do (set REVISION=%%i)
for /f "delims=" %%i in ('git rev-parse --short HEAD') do (set REVISION_HASH=%%i)

if "%REVISION%" == "" (
	set REVISION=0
) 

(echo #define VER_REVISION %REVISION%  && echo #define VER_REVISION_HASH %REVISION_HASH%) > %2
