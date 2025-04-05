# snippetbox

To config tls certificate
```bash
mkdir tls
cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
cd ..
```
To start the database
```bash
docker compose up
```
To run the application
```bash
go run ./cmd/web [>>/tmp/info.log 2>>/tmp/error.log]
```

## test
```bash
go test -covermode=count -coverprofile=/tmp/profile.out ./...
go tool cover -html=/tmp/profile.out
```
