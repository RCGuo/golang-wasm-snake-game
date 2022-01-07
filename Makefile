build_wasm:
	GOOS=js GOARCH=wasm go build -o ./site/snake.wasm .

serve:
	go run ./site/server.go