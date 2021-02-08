go tool compile -I kgen blissWebApp.go
go tool link -o blissLib.exe -L kgen blissWebApp.o
pause
