go tool compile -I kgen blissWebApp.go
go tool link -o blissAndLib.exe -L kgen blissWebApp.o
pause
