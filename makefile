ifeq ($(OS),Windows_NT)
	BINARY = start.exe
	CLEAN_CMD = del
else
	BINARY = start
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all:
	go build undra.go

.PHONY: clean
clean:
	$(BINARY) $(BINARY)