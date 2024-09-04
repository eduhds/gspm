go run . add go-task/task

Expand-Archive -LiteralPath $env:USERPROFILE\Downloads\task_windows_amd64.zip -DestinationPath $env:USERPROFILE\Downloads\task

copy $env:USERPROFILE\Downloads\task\task.exe .
