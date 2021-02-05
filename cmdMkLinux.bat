SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=linux
go build -o blissApp blissWebApp.go
pause
