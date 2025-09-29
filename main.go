package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	libvirt "github.com/libvirt/libvirt-go"
)

func main() {
	// 명령줄 인수 정의
	var (
		user = flag.String("user", "", "SSH 사용자명 (필수)")
		host = flag.String("host", "", "원격 호스트 IP 또는 호스트명 (필수)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법: %s -user <사용자명> -host <호스트>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "옵션:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n예시: %s -user suser -host 192.168.6.66\n", os.Args[0])
	}
	flag.Parse()

	// 필수 매개변수 검증
	if *user == "" || *host == "" {
		fmt.Fprintf(os.Stderr, "오류: user와 host 매개변수가 모두 필요합니다.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// 1) RO로 원격 접속 (socket 경로를 /run, /var/run 둘 다 시도)
	base := fmt.Sprintf("qemu+ssh://%s@%s/system", *user, *host)
	var conn *libvirt.Connect
	var err error
	for _, sock := range []string{"/run/libvirt/libvirt-sock-ro", "/var/run/libvirt/libvirt-sock-ro"} {
		uri := fmt.Sprintf("%s?socket=%s", base, sock)
		conn, err = libvirt.NewConnectReadOnly(uri)
		if err == nil {
			break
		}
	}
	if conn == nil || err != nil {
		log.Fatalf("RO 접속 실패: %v", err)
	}
	defer conn.Close()

	// 2) 도메인 나열
	doms, err := conn.ListAllDomains(0)
	if err != nil {
		log.Fatalf("도메인 조회 실패: %v", err)
	}
	for _, dom := range doms {
		printDomain(&dom)
		dom.Free() // libvirt 객체 해제
	}
}

// 도메인 1개 출력
func printDomain(dom *libvirt.Domain) {
	name, _ := dom.GetName()

	// 상태/요약 정보 (GetInfo는 RO에서 안전)
	info, err := dom.GetInfo()
	if err != nil {
		log.Printf("[%s] GetInfo 실패: %v", name, err)
		return
	}
	fmt.Printf("[%s] state=%s  vcpu=%d  mem=%dKiB  cputime=%dns\n",
		name, stateString(info.State), info.NrVirtCpu, info.Memory, info.CpuTime)

	// XML에서 디스크/인터페이스 타겟명, MAC 추출
	xmlStr, err := dom.GetXMLDesc(0)
	if err != nil {
		log.Printf("  XMLDesc 오류: %v", err)
		return
	}
	disks, ifaces, macs := extractFromXML(xmlStr)

	// 디스크 통계
	for _, dev := range disks {
		bs, err := dom.BlockStats(dev)
		if err != nil {
			// cdrom/드라이버 미지원 등은 무시
			continue
		}
		fmt.Printf("  disk %-8s rd=%dB wr=%dB errs=%d\n", dev, bs.RdBytes, bs.WrBytes, bs.Errs)
	}

	// NIC 통계
	for _, dev := range ifaces {
		is, err := dom.InterfaceStats(dev)
		if err != nil {
			continue
		}
		fmt.Printf("  nic  %-8s rx=%dB tx=%dB (rxPkts=%d txPkts=%d)\n",
			dev, is.RxBytes, is.TxBytes, is.RxPackets, is.TxPackets)
	}

	// (옵션) MAC만 출력 – IP가 필요하면 DHCP leases나 외부 소스와 매칭
	if len(macs) > 0 {
		fmt.Printf("  macs: %s\n", strings.Join(macs, ", "))
	}
}

// ===== XML 파서(encoding/xml 토큰 스캔로 가볍게) =====

func extractFromXML(x string) (disks, ifaces, macs []string) {
	dec := xml.NewDecoder(strings.NewReader(x))
	inDisk := false
	diskIsActual := false
	inIface := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "disk":
				inDisk = true
				diskIsActual = hasAttr(t, "device", "disk")
			case "target":
				// 디스크의 target dev
				if inDisk && diskIsActual {
					if dev, ok := getAttr(t, "dev"); ok {
						disks = append(disks, dev)
					}
				}
				// 인터페이스의 target dev
				if inIface {
					if dev, ok := getAttr(t, "dev"); ok {
						ifaces = append(ifaces, dev)
					}
				}
			case "interface":
				inIface = true
			case "mac":
				if addr, ok := getAttr(t, "address"); ok {
					macs = append(macs, addr)
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "disk":
				inDisk = false
				diskIsActual = false
			case "interface":
				inIface = false
			}
		}
	}
	return
}

func hasAttr(se xml.StartElement, key, val string) bool {
	for _, a := range se.Attr {
		if a.Name.Local == key && a.Value == val {
			return true
		}
	}
	return false
}
func getAttr(se xml.StartElement, key string) (string, bool) {
	for _, a := range se.Attr {
		if a.Name.Local == key {
			return a.Value, true
		}
	}
	return "", false
}

// ===== 상태 문자열 매핑 =====

func stateString(s libvirt.DomainState) string {
	switch s {
	case libvirt.DOMAIN_NOSTATE:
		return "nostate"
	case libvirt.DOMAIN_RUNNING:
		return "running"
	case libvirt.DOMAIN_BLOCKED:
		return "blocked"
	case libvirt.DOMAIN_PAUSED:
		return "paused"
	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdown"
	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"
	case libvirt.DOMAIN_CRASHED:
		return "crashed"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "pmsuspended"
	default:
		return fmt.Sprintf("%d", s)
	}
}
