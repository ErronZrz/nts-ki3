package congrat2

import (
	"active/utils"
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"runtime"
	"testing"
	"time"
)

// 好吧，我彻底放弃了，TTL 直接用默认值算了，服气，Go 开发人员到底是有多懒，连个函数都不实现

func getTTL(conn *net.UDPConn) (int, error) {
	// 获取原始文件描述符
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return 0, err
	}

	var ttl int
	var syscallErr error
	errCtrl := rawConn.Control(func(fd uintptr) {
		// Windows 接收 TTL 的特殊选项
		const IP_TTL = 0x21 // IPPROTO_IP=0 + 21=IP_RECVTTL
		syscallErr = windows.SetsockoptInt(
			windows.Handle(fd),
			windows.IPPROTO_IP,
			windows.IP_TTL,
			5,
		)
		if syscallErr != nil {
			return
		}

		// 获取 TTL 值
		ttl, syscallErr = windows.GetsockoptInt(
			windows.Handle(fd),
			windows.IPPROTO_IP,
			windows.IP_TTL,
		)
	})

	if errCtrl != nil {
		return 0, errCtrl
	}
	return ttl, syscallErr
}

func TestSomething(t *testing.T) {
	serverIP := "200.160.7.197" // 替换为你的NTP服务器IP
	ntpRequest := utils.VariableData()

	addr, _ := net.ResolveUDPAddr("udp", serverIP+":123")
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	// 设置超时
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// 发送请求
	conn.Write(ntpRequest)

	// 接收响应
	buffer := make([]byte, 512)
	_, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	// 获取 TTL（需在接收数据后调用）
	ttl, err := getTTL(conn)
	if err != nil {
		if runtime.GOOS == "windows" {
			fmt.Println("注意：Windows 获取 TTL 需要管理员权限")
		}
		panic(err)
	}

	fmt.Printf("Received TTL: %d\n", ttl)
}
