package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

func HumanBytes(n int64) string {
	return humanize.Bytes(uint64(n))
}
func UnmarshalConfigFromFile(filePath string, config interface{}) error {
	data, err := GetFileContent(filePath)
	if err != nil {
		return err
	}
	if strings.HasSuffix(filePath, ".json") {
		return json.Unmarshal(data, config)
	}
	if strings.HasSuffix(filePath, ".yml") || strings.HasSuffix(filePath, ".yaml") {
		return yaml.Unmarshal(data, config)
	}
	return fmt.Errorf("unknown extension of file: %s", filePath)
}
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// getDefaultDevice 查找具有 IP 地址且状态为 UP 的网卡
func GetDefaultDevice() (string, error) {
	// 获取所有网卡
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("Error getting interfaces: %v", err)
	}

	// 遍历网卡并检查每个网卡的状态
	for _, iface := range interfaces {
		// 获取网卡的标志
		flags := iface.Flags

		// 检查网卡是否处于 UP 状态
		if flags&net.FlagUp == 0 {
			continue
		}

		// 获取网卡的地址信息
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		// 检查网卡是否有有效的 IP 地址（包括 IPv4 和 IPv6）
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipnet.IP

			// 排除回环地址
			if ip.IsLoopback() {
				continue
			}

			// 检查是否为有效的 IP 地址
			if ip.To4() != nil || ip.To16() != nil {
				return iface.Name, nil
			}
		}
	}

	return "", fmt.Errorf("No suitable device found")
}

// 监控特定端口的上传流量
type TrafficMonitor struct {
	device       string
	port         string
	handle       *pcap.Handle
	packetChan   chan gopacket.Packet
	TotalTxBytes int64
	mu           sync.Mutex
}

// NewTrafficMonitor 创建一个新的上传流量监控器
func NewTrafficMonitor(device string, port string) (*TrafficMonitor, error) {
	if device == "" {
		var err error
		device, err = GetDefaultDevice()
		println(device)
		println("device_id")
		if err != nil {
			return nil, fmt.Errorf("Error selecting default device: %v", err)
		}
	}

	handle, err := pcap.OpenLive(device, 160000, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("Error opening device %s: %v", device, err)
	}

	// 设置 BPF 过滤器，过滤特定端口的流量
	filter := fmt.Sprintf("tcp and port %s", port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		handle.Close()
		return nil, fmt.Errorf("Error setting BPF filter: %v", err)
	}

	monitor := &TrafficMonitor{
		device:     device,
		port:       port,
		handle:     handle,
		packetChan: make(chan gopacket.Packet, 100),
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		for packet := range packetSource.Packets() {
			monitor.packetChan <- packet
		}
	}()

	go monitor.processPackets()

	return monitor, nil
}

// processPackets 处理数据包并累积上传流量
func (tm *TrafficMonitor) processPackets() {
	for packet := range tm.packetChan {
		if packet.NetworkLayer() != nil && packet.TransportLayer() != nil {
			tm.mu.Lock()
			tm.TotalTxBytes += int64(len(packet.Data()))
			tm.mu.Unlock()
		}
	}
}

// GetUploadTraffic 返回从上次调用以来的上传流量（字节）
func (tm *TrafficMonitor) GetUploadTraffic() int64 {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	uploadBytes := tm.TotalTxBytes
	tm.TotalTxBytes = 0 // 重置计数器，获取到的就是自上次调用以来的流量
	return uploadBytes
}

// Close 关闭监控器
func (tm *TrafficMonitor) Close() {
	if tm.handle != nil {
		tm.handle.Close()
	}
}
