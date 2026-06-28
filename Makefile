BINARY_NAME=hashto

ifeq ($(OS),Windows_NT)
    TARGET_OS=windows
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        TARGET_ARCH=amd64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            TARGET_ARCH=amd64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),ARM64)
            TARGET_ARCH=arm64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            TARGET_ARCH=386
        endif
    endif
    EXE_EXT=.exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        TARGET_OS=linux
    endif
    ifeq ($(UNAME_S),Darwin)
        TARGET_OS=darwin
    endif

    UNAME_M := $(shell uname -m)
    ifeq ($(UNAME_M),x86_64)
        TARGET_ARCH=amd64
    endif
    ifeq ($(UNAME_M),arm64)
        TARGET_ARCH=arm64
    endif
    ifeq ($(UNAME_M),aarch64)
        TARGET_ARCH=arm64
    endif
    EXE_EXT=
endif

OUTPUT_NAME=$(BINARY_NAME)-$(TARGET_OS)-$(TARGET_ARCH)$(EXE_EXT)

all: build

build:
	@echo "Detecting OS and Architecture..."
	@echo "Target OS: $(TARGET_OS)"
	@echo "Target Arch: $(TARGET_ARCH)"
	@echo "Building binary: $(OUTPUT_NAME)..."
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $(OUTPUT_NAME) main.go compiler.go
	@echo "Build complete!"

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)-*
