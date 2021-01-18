ifeq ($(OS),Windows_NT)
	BINARY = undrastart.exe
	CLEAN_CMD = del
else
	SET_GOPATH = GOPATH=$(GOPATH)
	BINARY = undrastart.out
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: default
default: all

.PHONY: all
all:
	$(SET_GOPATH) go build -o $(BINARY) undrastart.go

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)