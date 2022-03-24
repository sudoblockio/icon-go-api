go get github.com/swaggo/swag/cmd/swag
go get github.com/arsmn/fiber-swagger/v2@v2.20.0
go get github.com/alecthomas/template
go get github.com/riferrei/srclient@v0.3.0
swag init --parseDependency -g api.go
