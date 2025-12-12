.PHONY: all build run clean generate

INPUT = input/data.txt

all: generate build run

generate:
	@echo "Generating input data..."
	@go run scripts/generate-input.go

build: build-go-stdlib build-go-coregex

build-go-stdlib:
	@echo "Building go-stdlib..."
	@cd go-stdlib && go build -ldflags "-s -w" -o ../bin/go-stdlib.exe .

build-go-coregex:
	@echo "Building go-coregex..."
	@cd go-coregex && go mod tidy && go build -ldflags "-s -w" -o ../bin/go-coregex.exe .

run: $(INPUT)
	@echo ""
	@echo "==============================================="
	@./bin/go-stdlib.exe $(INPUT)
	@echo ""
	@echo "==============================================="
	@./bin/go-coregex.exe $(INPUT)
	@echo ""

clean:
	@rm -rf bin/*.exe input/data.txt

$(INPUT):
	@$(MAKE) generate
