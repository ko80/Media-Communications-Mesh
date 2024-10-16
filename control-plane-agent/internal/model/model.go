package model

import "time"

type MediaProxy struct {
	ID     string            `json:"id,omitempty"`
	Config *MediaProxyConfig `json:"config,omitempty"`
	Status *MediaProxyStatus `json:"status,omitempty"`
}

type MediaProxyConfig struct {
	IPAddr string `json:"ipAddr"`
}

type MediaProxyStatus struct {
	Healthy       bool      `json:"healthy"`
	StartedAt     time.Time `json:"startedAt"`
	ConnectionNum int       `json:"connectionNum"`
}

type Connection struct {
	ID     string            `json:"id,omitempty"`
	Config *ConnectionConfig `json:"config,omitempty"`
	Status *ConnectionStatus `json:"status,omitempty"`
}

type ConnectionConfig struct {
	Kind        string      `json:"kind"`
	ConnType    string      `json:"connType"`
	Conn        interface{} `json:"conn"`
	PayloadType string      `json:"payloadType"`
	Payload     interface{} `json:"payyload"`
	BufferSize  uint64      `json:"bufferSize"`
}

type ConnectionST2110 struct {
	RemoteIPAddr string `json:"remoteIpAddr"`
	RemotePort   string `json:"remotePort"`
}

// TBD ConnectionMemif
// TBD ConnectionRDMA

type PayloadVideo struct {
	Width       uint32  `json:"width"`
	Height      uint32  `json:"height"`
	FPS         float32 `json:"fps"`
	PixelFormat string  `json:"pixelFormat"`
}

// TBD PayloadAudio
// TBD PayloadAncillary

type ConnectionStatus struct {
	CreatedAt   time.Time `json:"createdAt"`
	BuffersSent int       `json:"buffersSent"`
}

type MultipointGroup struct {
	ID     string                 `json:"id,omitempty"`
	Config *MultipointGroupConfig `json:"config,omitempty"`
	Status *MultipointGroupStatus `json:"status,omitempty"`
}

type MultipointGroupConfig struct {
	IPAddr string `json:"ipAddr"`
}

type MultipointGroupStatus struct {
	CreatedAt time.Time `json:"createdAt"`
}

type Bridge struct {
	ID     string        `json:"id,omitempty"`
	Config *BridgeConfig `json:"config,omitempty"`
	Status *BridgeStatus `json:"status,omitempty"`
}

type BridgeConfig struct {
	IPAddr string `json:"ipAddr"`
}

type BridgeStatus struct {
	CreatedAt time.Time `json:"createdAt"`
}
