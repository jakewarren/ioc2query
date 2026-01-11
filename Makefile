.PHONY: build wasm clean test run-web

BINARY_NAME=ioc2query
WASM_OUT=web/main.wasm
WASM_EXEC_JS=web/wasm_exec.js

build:
	go build -o $(BINARY_NAME) ./cmd/ioc2query

wasm: $(WASM_OUT) $(WASM_EXEC_JS)

$(WASM_OUT): cmd/ioc2query/main.go pkg/extraction/extractor.go pkg/backends/s1ql/s1ql.go
	@echo "Building WASM binary..."
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) ./cmd/ioc2query

$(WASM_EXEC_JS):
	@echo "Copying wasm_exec.js..."
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_EXEC_JS) || cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_EXEC_JS)

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f $(WASM_OUT)

run-web: wasm
	@echo "Starting web server at http://localhost:8080"
	python3 -m http.server -d web 8080
