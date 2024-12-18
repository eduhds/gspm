rd /s /q build
md build
md build\windows
md build\windows\amd64

rd /s /q dist
md dist

del *.syso
go-winres make --in windows/winres.json

set CGO_ENABLED=0
go build -ldflags "-s -w" -o build\windows\amd64\gspm.exe .

copy README.md build\windows\amd64
copy LICENSE.txt build\windows\amd64

tar -C build\windows\amd64 -a -c -f dist\gspm-windows-amd64.zip gspm.exe README.md LICENSE.txt
start .\windows\setup.iss
