ifeq ($(OS),Windows_NT)
	BINARY = undrastart.exe
	CLEAN_CMD = del
else
	BINARY = undrastart
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: default
default: all

.PHONY: all
all:
	go build undrastart.go

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)