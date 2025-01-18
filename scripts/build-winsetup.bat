@echo off

rem rd /s /q build
rem md build
rem md build\windows
rem md build\windows\amd64

rem rd /s /q dist
rem md dist

rem del *.syso
rem go-winres make --in windows/winres.json

rem set CGO_ENABLED=0
rem go build -ldflags "-s -w" -o build\windows\amd64\gspm.exe .

rem copy README.md build\windows\amd64
rem copy LICENSE.txt build\windows\amd64

rem tar -C build\windows\amd64 -a -c -f dist\gspm-windows-amd64.zip gspm.exe README.md LICENSE.txt

start .\windows\setup.iss
