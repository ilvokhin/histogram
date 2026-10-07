BIN=bin

.PHONY: all clean test

all: $(BIN)/histogram

$(BIN)/histogram:
	go build -o $(BIN)/histogram

test:
	go test

clean:
	rm -rf $(BIN)
