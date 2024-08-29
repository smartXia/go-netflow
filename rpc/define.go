package rpc

import (
	"encoding/json"
)

type DeviceInfo struct {
	Cpus      []CpuInfo     `json:"cpus" yaml:"cpus"` // CPU信息
	Disks     []DiskInfo    `json:"disks" yaml:"disks"`
	Hostname  string        `json:"hostname" yaml:"hostname"` // 主机名
	Memories  []MemoryInfo  `json:"memories" yaml:"memories"`
	Model     string        `json:"model" yaml:"model"` // 设备型号
	NetCards  []NetCardInfo `json:"netCards" yaml:"netCards"`
	OSVersion string        `json:"osVersion" yaml:"osVersion"` // 镜像版本
	Sn        string        `json:"sn" yaml:"sn"`               // 设备sn
}

type CpuInfo struct {
	CacheSize  int64       `json:"cacheSize" yaml:"cacheSize"`
	Cores      int64       `json:"cores" yaml:"cores"`
	CpuNO      int64       `json:"cpu" yaml:"cpu"` // CPU编号
	Family     string      `json:"family" yaml:"family"`
	Flags      interface{} `json:"flags" yaml:"flags"`
	Ghz        int64       `json:"ghz" yaml:"ghz"`
	Mhz        int64       `json:"mhz" yaml:"mhz"` // 主频
	Microcode  string      `json:"microcode" yaml:"microcode"`
	Model      string      `json:"model" yaml:"model"`
	ModelName  string      `json:"modelName" yaml:"modelName"` // 型号
	Name       string      `json:"name" yaml:"name"`
	PhysicalID string      `json:"physicalId" yaml:"physicalId"`
	Stepping   int64       `json:"stepping" yaml:"stepping"`
	VendorID   string      `json:"vendorId" yaml:"vendorId"` // 厂商
}

type DiskInfo struct {
	ID           string `json:"id"`           //硬盘名称
	Size         uint64 `json:"size"`         //磁盘容量
	Media        string `json:"media"`        //类型
	SerialNumber string `json:"serialNumber"` //序列号
}

type MemoryInfo struct {
	Mhz      int64  `json:"mhz" yaml:"mhz"`                   // 速度
	SerialNO string `json:"serialNumber" yaml:"serialNumber"` // 序列号
	Size     int64  `json:"size" yaml:"size"`                 // 大小
	Slot     string `json:"slot" yaml:"slot"`                 // 插槽
	Type     string `json:"type" yaml:"type"`                 // 型号(DDR4)
}

// 网卡列表
type NetCardInfo struct {
	DialingInfo []DialingInfo `json:"dialingInfo" yaml:"dialingInfo"` // 拨号信息.
	IP          string        `json:"ip" yaml:"ip"`                   // ip.
	IsManager   bool          `json:"isManager" yaml:"isManager"`     // 是否管理网卡
	IsValid     bool          `json:"isValid" yaml:"isValid"`         // 是否有效
	Name        string        `json:"name" yaml:"name"`               // 网卡名称
	Speed       int64         `json:"speed" yaml:"speed"`             // 速率 单位M.
}

type DialingInfo struct {
	Account      string `json:"account" yaml:"account"`
	Connected    bool   `json:"connected" yaml:"connected"`
	EnableV6     int64  `json:"enableV6" yaml:"enableV6"`
	FailedReason string `json:"failedReason" yaml:"failedReason"`
	Gateway      string `json:"gateway" yaml:"gateway"`
	IP           string `json:"ip" yaml:"ip"`
	PppoeName    string `json:"pppoeName" yaml:"pppoeName"`
	Port         int    `json:"port" yaml:"port"`
	MAC          string `json:"mac" yaml:"mac"`
	Mask         string `json:"mask" yaml:"mask"`
	Number       int64  `json:"number" yaml:"number"`
	Password     string `json:"password" yaml:"password"`
	PppoeAC      string `json:"pppoeAC" yaml:"pppoeAC"`
	PppoeService string `json:"pppoeService" yaml:"pppoeService"`
	Successed    bool   `json:"successed" yaml:"successed"`
	V4Proto      string `json:"v4Proto" yaml:"v4Proto"`
	V6Connected  bool   `json:"v6Connected" yaml:"v6Connected"`
	V6Gateway    string `json:"v6Gateway" yaml:"v6Gateway"`
	V6IP         string `json:"v6IP" yaml:"v6IP"`
	V6MaskLen    int64  `json:"v6MaskLen" yaml:"v6MaskLen"`
	V6Proto      string `json:"v6Proto" yaml:"v6Proto"`
	V6Successed  bool   `json:"v6Successed" yaml:"v6Successed"`
	VLANID       string `json:"vlanId" yaml:"vlanId"`
}

type RegisterRsp struct {
	CodeURL      string `json:"codeURL" yaml:"codeURL"`   // 二维码地址
	DeviceID     string `json:"deviceID" yaml:"deviceID"` // 设备ID
	FrpcPort     int    `json:"frpcPort" yaml:"frpcPort"`
	FrpsHost     string `json:"frpsHost" yaml:"frpsHost"`
	FrpsPort     int    `json:"frpsPort" yaml:"frpsPort"`
	InfluxDBHost string `json:"influxDBHost" yaml:"influxDBHost"`
	InfluxDBPort int    `json:"influxDBPort" yaml:"influxDBPort"`
}

type HeartbeatRsp struct {
	DiskTestTask      DiskTestTask    `json:"diskTestTask" yaml:"diskTestTask"`
	Netcards          Netcards        `json:"netcards" yaml:"netcards"`
	NetworkTestTask   NetworkTestTask `json:"networkTestTask" yaml:"networkTestTask"`
	NextHeartBeatTime int64           `json:"nextHeartBeatTime" yaml:"nextHeartBeatTime"`
	Upgrade           Upgrade         `json:"upgrade" yaml:"upgrade"`
	Deployment        Deployment      `json:"deploy" yaml:"deploy"`
}

type MonitorInfo struct {
	CPUUsage      float64 `json:"cpuUsage"`
	DiskUsage     float64 `json:"diskUsage"`
	DownBandwidth float64 `json:"downBandwidth"`
	MemUsage      float64 `json:"memUsage"`
	Timestamp     int64   `json:"timestamp"`
	UpBandwidth   float64 `json:"upBandwidth"`
}

type Netcards struct {
	ModifyTime  int64         `json:"modifyTime" yaml:"modifyTime"`
	DialingType string        `json:"dialingType" yaml:"dialingType"`
	Netcards    []NetCardInfo `json:"data" yaml:"data"`
}

type NetworkTestTask struct {
	IsValid         bool              `json:"isValid" yaml:"isValid"`
	Time            int64             `json:"time" yaml:"time"`
	NetworkTestInfo []NetworkTestInfo `json:"networkTestInfo" yaml:"networkTestInfo"`
}

type NetworkTestInfo struct {
	IP    string  `json:"ip" yaml:"ip"`
	Limit *string `json:"limit" yaml:"limit"`
	Port  int     `json:"port" yaml:"port"`
}

type DiskTestTask struct {
	IsValid bool  `json:"isValid" yaml:"isValid"`
	Time    int64 `json:"time" yaml:"time"`
}

type Upgrade struct {
	IsValid bool   `json:"isValid" yaml:"isValid"`
	MD5     string `json:"md5" yaml:"md5"`
	URL     string `json:"url" yaml:"url"`
}

type Deployment struct {
	Action        string           `json:"action" yaml:"action"`
	Bandwidth     int              `json:"bandwidth" yaml:"bandwidth"` //部署的单条带宽大小
	CMD           string           `json:"cmd" yaml:"cmd"`
	DeploymentId  string           `json:"deploymentId" yaml:"deploymentId"`
	Desc          string           `json:"desc" yaml:"desc"`
	DeviceInfo    DeployDeviceInfo `json:"deviceInfo" yaml:"deviceInfo"`
	From          string           `json:"from" yaml:"from"`
	IsOutProvince int              `json:"isOutProvince" yaml:"isOutProvince"`
}

type DeployDeviceInfo struct {
	BizType         int    `json:"biz_type" yaml:"biz_type"`
	DeviceId        string `json:"device_id" yaml:"device_id"`
	SingleBandwidth uint64 `json:"single_bandwidth" yaml:"single_bandwidth"`
	SubBiztype      int    `json:"sub_biz_type" yaml:"sub_biz_type"`
}

type DeploymentStatusResult struct {
	DeploymentId string `json:"deploymentId" yaml:"deploymentId"`
	Status       int    `json:"status" yaml:"status"`
	FailReason   string `json:"failReason" yaml:"failReason"`
}

type NetworkTestResult struct {
	Type                   string        `json:"type" yaml:"type"`                                     // 压测类型，可选(丢包测试,极限测试
	MaxTestBandwidth       string        `json:"maxTestBandwidth" yaml:"maxTestBandwidth"`             // 最大带宽，单位：bps.
	AvgTestBandwidth       string        `json:"avgTestBandwidth" yaml:"avgTestBandwidth"`             // 平均带宽，单位：bps.
	UpBandwidthTestTime    string        `json:"upBandwidthTestTime" yaml:"upBandwidthTestTime"`       // 测试时间，UNIX时间戳.
	Duration               string        `json:"duration" yaml:"duration"`                             // 用时，单位：s.
	LineQuality            []LineQuality `json:"lineQuality" yaml:"lineQuality"`                       // 线路信息.
	TCPRetransmissionRatio int64         `json:"tcpRetransmissionRatio" yaml:"tcpRetransmissionRatio"` // TCP 重传率.
}

type LineQuality struct {
	Name           string `json:"name" yaml:"name"` // 线路网卡名.
	IP             string `json:"ip" yaml:"ip"`
	AvgUpBandwidth string `json:"avgUpBandwidth" yaml:"avgUpBandwidth"` // 平均速率，单位 bps.
	NetworkDelay   int    `json:"networkDelay" yaml:"networkDelay"`     // 网络延时，单位：ms.
	PacketLossRate int    `json:"packetLossRate" yaml:"packetLossRate"` // 丢包率， 0.05 表示丢包率为 5%.
}

type DiskResults struct {
	Type        string        `json:"type" yaml:"type"` // 测试类型，(随机读/写）
	DiskQuality []DiskQuality `json:"diskQuality" yaml:"diskQuality"`
}

type DiskQuality struct {
	Id        string `json:"id" yaml:"id"`               // 物理磁盘id /dev/sda
	Iops      int    `json:"iops" yaml:"iops"`           //  读IOPS
	Bandwidth int    `json:"bandwidth" yaml:"bandwidth"` // 带宽，最大吞吐量（MB/s）
}

type ServerError struct {
	Code    int    `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
	EventID string `json:"event_id" yaml:"event_id"`
}

func (e *ServerError) Error() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

type HttpError struct {
	Code    int    `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}

func (e *HttpError) Error() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
