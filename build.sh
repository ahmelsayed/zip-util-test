# GOOS=linux
# GOOS=windows

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o go-zip main.go
