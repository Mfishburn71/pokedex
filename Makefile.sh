#GOOS=js GOARCH=wasm go build -o main.wasm .

.PHONY: build serve clean

build:
	GOOS=js GOARCH=wasm go build -o docs/main.wasm .
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" docs/

serve:
	python3 -m http.server 8080 --directory docs

clean:
	rm -f docs/main.wasm docs/wasm_exec.js

build-native:
	go build -o pokedex .