ifeq ($(OS),Windows_NT)
	BINARY = undra_start.exe
	CLEAN_CMD = del
else
	BINARY = undra_start
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../

.PHONY: all
all:
	go build undra_start.go

.PHONY: clean
clean:
	$(CLEAN_CMD) $(BINARY)