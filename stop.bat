@echo off

set dir=%~dp0
cd /D %dir%
mergernew.exe -cmd stop -d true

pause
