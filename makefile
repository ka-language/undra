ifeq ($(OS),Windows_NT)
	BINARY = undrastart.exe
	CLEAN_CMD = del
else
	BINARY = undrastart
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: go.mod
	go get -u
	go build undrastart.go

go.mod:
	go mod init

.PHONY: clean
clean:
	-go mod tidy
	-$(CLEAN_CMD) $(BINARY)