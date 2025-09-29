# ====== Config ======
BIN       ?= lvmon          # 출력 바이너리 이름
GO        ?= go
GOFLAGS   ?=
LDFLAGS   ?=
CGO       ?= 1              # libvirt-go는 CGO 필요

# ====== Phony targets ======
.PHONY: all build run tidy fmt vet clean check

all: build

# libvirt dev 패키지/헤더 확인(권장)
check:
	@command -v pkg-config >/dev/null || (echo "ERR: pkg-config not found" && exit 1)
	@pkg-config --exists libvirt || (echo "ERR: libvirt development package not found" && exit 1)
	@echo "libvirt version: $$(pkg-config --modversion libvirt)"

# 빌드(현재 디렉터리의 모듈을 main으로 빌드)
build: check go.mod go.sum
	CGO_ENABLED=$(CGO) $(GO) build $(GOFLAGS) -o $(BIN) -ldflags '$(LDFLAGS)' .

# 실행(필요시)
run: build
	./$(BIN)

go.mod:
	$(GO) mod init ntels.com/libvirt-vm-mon

# 모듈 정리/검사(선택)
go.sum:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet: $(GO) vet ./...

# 바이너리 정리
clean:
	@rm -f $(BIN)
