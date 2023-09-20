go get -d github.com/swaggo/swag/cmd/swag
go get -d github.com/arsmn/fiber-swagger/v2@v2.31.1
go get -d github.com/alecthomas/template
go get -d github.com/riferrei/srclient@v0.3.0
swag init --parseDependency -g api.go

# Note if you get the error
# cannot find all dependencies, unable to resolve root package
# See https://github.com/swaggo/swag/issues/909#issuecomment-834107747
# Using wrong version. Should be same as go.mod
