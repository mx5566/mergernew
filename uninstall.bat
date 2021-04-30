@echo off

set dir=%~dp0
cd /D %dir%

mergernew.exe -cmd uninstall -d true

pause