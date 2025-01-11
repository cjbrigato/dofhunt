go-winres make
go build -o dofhunt-win64.exe -ldflags "-s -w -H=windowsgui -extldflags=-static" .
