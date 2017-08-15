//go:generate sh -c "GOOS=linux GOARCH=amd64 go build -o up-proxy ../../../cmd/up-proxy/main.go"
//go:generate go-bindata -modtime 0 -pkg bin -o bin_assets.go .

package bin
