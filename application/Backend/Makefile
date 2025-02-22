TARGET_OS = linux
TARGET_ARCH = amd64

#Mac bianary adjustment
UNAME_S = $(shell uname -s)
ifeq ($(UNAME_S), Darwin)
	TARGET_ARCH = arm64
endif

OUTPUT_NAME = beckend-$(TARGET_OS)-$(TARGET_ARCH)

vendor:
	GOWORK=off go mod vendor

run:
	GOWORK=off go run main.go

build:
	GOWORK=off GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $(OUTPUT_NAME) main.go

build-run: build
	./$(OUTPUT_NAME)

#Unsafe version definelty need to fix!
allow-port:
	sudo ufw allow 8081/tcp

remove-port:
	sudo ufw delete allow 8081/tcp

clean:
	rm -f backend-*
