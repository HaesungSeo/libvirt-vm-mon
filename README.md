# libvirt-vm-mon

libvirt를 통해 원격 서버의 가상머신 상태와 리소스 사용량을 모니터링하는 Go 프로그램입니다.

## 기능

- 원격 libvirt 서버에 SSH를 통해 읽기 전용으로 접속
- 가상머신 목록 조회 및 상태 확인
- 각 VM의 CPU, 메모리, 디스크, 네트워크 사용량 통계 출력
- MAC 주소 정보 표시

## 빌드 방법

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
./libvirt-vm-mon -user suser -host 192.168.6.66

# 호스트명으로 접속
./libvirt-vm-mon -user admin -host vm-host.example.com

# 도움말 보기
./libvirt-vm-mon -h
```

## 출력 예시

```
[vm1] state=running  vcpu=2  mem=2097152KiB  cputime=12345678ns
  disk vda      rd=1048576B wr=524288B errs=0
  disk vdb      rd=2097152B wr=1048576B errs=0
  nic  vnet0    rx=1024B tx=2048B (rxPkts=10 txPkts=20)
  macs: 52:54:00:12:34:56

[vm2] state=shutoff  vcpu=1  mem=1048576KiB  cputime=0ns
  macs: 52:54:00:78:9a:bc
```

## 요구사항

- Go 1.18 이상
- libvirt-go 패키지 (`github.com/libvirt/libvirt-go`)
- 원격 서버에 SSH 키 기반 인증 설정
- 원격 서버에 libvirt 설치 및 실행 중

## 주의사항

- 이 프로그램은 읽기 전용 모드로 작동하므로 VM을 제어하거나 설정을 변경할 수 없습니다
- SSH 키 기반 인증이 설정되어 있어야 합니다
- 원격 서버의 libvirt 소켓 경로는 `/run/libvirt/libvirt-sock-ro` 또는 `/var/run/libvirt/libvirt-sock-ro`를 시도합니다
