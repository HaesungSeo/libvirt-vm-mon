# libvirt-vm-mon

libvirt를 통해 원격 서버의 가상머신 상태와 리소스 사용량을 모니터링하는 Go 프로그램입니다.

## 기능

- 원격 libvirt 서버에 SSH를 통해 읽기 전용으로 접속
- 가상머신 목록 조회 및 상태 확인
- 각 VM의 CPU, 메모리, 디스크, 네트워크 사용량 통계 출력
- MAC 주소 정보 표시

## 설치 및 빌드

### 1. 프로젝트 다운로드

#### Git을 사용하여 클론:
```bash
git clone https://github.com/HaesungSeo/libvirt-vm-mon.git
cd libvirt-vm-mon
```

#### 또는 소스 코드를 직접 다운로드한 경우:
```bash
# 압축 해제 후 디렉토리로 이동
cd libvirt-vm-mon
```

### 2. 의존성 확인 및 설치

Go 모듈이 이미 설정되어 있으므로 의존성은 자동으로 다운로드됩니다:

```bash
# 의존성 확인
go mod tidy

# 의존성 다운로드 (필요시)
go mod download
```

### 3. 빌드 방법

```bash
go build -o libvirt-vm-mon main.go
```

또는 Makefile을 사용하는 경우:
```bash
make
```

## 사용법

```bash
./libvirt-vm-mon -user <사용자명> -host <호스트>
```

### 매개변수

- `-user`: SSH 접속에 사용할 사용자명 (필수)
- `-host`: 원격 호스트의 IP 주소 또는 호스트명 (필수)
- `-h`: 도움말 표시

### 사용 예시

```bash
# 기본 사용법
./libvirt-vm-mon -user myuser -host 192.168.1.100

# 또는 make로 빌드한 경우
./lvmon -user myuser -host 192.168.1.100

# 호스트명으로 접속
./libvirt-vm-mon -user admin -host vm-host.example.com

# 도움말 보기
./libvirt-vm-mon -h
```

### 실행 결과 예시

원격 서버에 `web-server`와 `database-server` 두 개의 VM이 실행 중인 경우:

```bash
$ ./lvmon -user myuser -host 192.168.1.100
[web-server] state=running  vcpu=4  mem=4194304KiB  cputime=1763967430000000ns
  disk vda      rd=323655595008B wr=230157305344B errs=0
  nic  web-server-eth0 rx=12314334870B tx=845075268B (rxPkts=56676469 txPkts=3532344)
  macs: 52:54:00:ab:cd:ef

[database-server] state=running  vcpu=16  mem=25165824KiB  cputime=32130911240000000ns
  disk vda      rd=131380239360B wr=10956509222912B errs=0
  nic  database-server-eth0 rx=681668217242B tx=7611180865198B (rxPkts=3170865327 txPkts=3033473839)
  macs: 52:54:00:12:34:56
```

#### 출력 정보 설명

각 VM에 대해 다음 정보가 표시됩니다:

- **VM 이름과 기본 상태**
  - `state`: VM 실행 상태 (running, shutoff, paused 등)
  - `vcpu`: 할당된 가상 CPU 개수
  - `mem`: 할당된 메모리 (KiB 단위)
  - `cputime`: 누적 CPU 사용 시간 (나노초 단위)

- **디스크 통계** (`disk` 항목)
  - `rd`: 읽은 바이트 수
  - `wr`: 쓴 바이트 수
  - `errs`: 디스크 오류 횟수

- **네트워크 인터페이스 통계** (`nic` 항목)
  - `rx`: 받은 바이트 수
  - `tx`: 보낸 바이트 수
  - `rxPkts`: 받은 패킷 수
  - `txPkts`: 보낸 패킷 수

- **MAC 주소**: VM의 네트워크 인터페이스 MAC 주소

## 요구사항

### 로컬 시스템 (빌드 및 실행 환경)
- Go 1.18 이상
- Git (프로젝트 다운로드용)
- libvirt 개발 라이브러리 (libvirt-dev 또는 libvirt-devel)

#### Ubuntu/Debian에서 libvirt 개발 라이브러리 설치:
```bash
sudo apt-get update
sudo apt-get install libvirt-dev pkg-config
```

#### CentOS/RHEL/Fedora에서 libvirt 개발 라이브러리 설치:
```bash
# CentOS/RHEL
sudo yum install libvirt-devel pkgconfig

# Fedora
sudo dnf install libvirt-devel pkgconfig
```

### 원격 서버 (모니터링 대상)
- libvirt 설치 및 실행 중
- SSH 서버 실행 중
- SSH 키 기반 인증 설정 (패스워드 없는 로그인)

## Go 의존성
- `github.com/libvirt/libvirt-go` v7.4.0+incompatible

## 주의사항

- 이 프로그램은 읽기 전용 모드로 작동하므로 VM을 제어하거나 설정을 변경할 수 없습니다
- SSH 키 기반 인증이 설정되어 있어야 합니다
- 원격 서버의 libvirt 소켓 경로는 `/run/libvirt/libvirt-sock-ro` 또는 `/var/run/libvirt/libvirt-sock-ro`를 시도합니다

## 문제 해결

### 빌드 오류 해결

#### "libvirt.h: No such file or directory" 오류
```bash
# Ubuntu/Debian
sudo apt-get install libvirt-dev pkg-config

# CentOS/RHEL/Fedora
sudo yum install libvirt-devel pkgconfig  # 또는 dnf
```

#### "pkg-config not found" 오류
```bash
# Ubuntu/Debian
sudo apt-get install pkg-config

# CentOS/RHEL/Fedora  
sudo yum install pkgconfig  # 또는 dnf
```

### 실행 시 오류 해결

#### "접속 실패" 오류
1. SSH 키 인증이 설정되어 있는지 확인:
   ```bash
   ssh user@host "echo 'SSH 연결 성공'"
   ```
2. 원격 서버에서 libvirt 서비스 상태 확인:
   ```bash
   ssh user@host "sudo systemctl status libvirtd"
   ```

#### "권한 거부" 오류
원격 서버에서 사용자가 libvirt 그룹에 속해 있는지 확인:
```bash
ssh user@host "groups"
# libvirt 그룹이 없다면:
ssh user@host "sudo usermod -aG libvirt $USER"
```
