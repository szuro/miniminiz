package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"net"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ZABBIX_FLAG       = 1
	COMPRESSION_FLAG  = 2
	LARGE_PACKET_FLAG = 4
)

var HEADER = []byte{'Z', 'B', 'X', 'D', ZABBIX_FLAG}

type ActiveItem struct {
	Key     string `json:"key"`
	Delay   int32  `json:"delay" yaml:"interval"`
	Logsize int32  `json:"lastlogsize"`
	MTime   int32  `json:"mtime"`
}

type ActiveItemValue struct {
	Host  string      `json:"host"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Clock int32       `json:"clock"`
	Ns    int32       `json:"ns"`
}

type AgentRequest struct {
	Host    string            `json:"host"`
	Request string            `json:"request"`
	Session string            `json:"session"`
	Data    []ActiveItemValue `json:"data"`
}

type ServerResponse struct {
	Response string       `json:"response"`
	Data     []ActiveItem `json:"data"`
	Info     string       `json:"info"`
}

type Server struct {
	connection string
	Monitoring MonitoringConfig
	Cache      chan ActiveItemValue
}

func NewServer(configFile string) (server *Server, err error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)

	if config.Server.cacheSize == 0 {
		config.Server.cacheSize = 100
	}

	server = &Server{}
	server.Cache = make(chan ActiveItemValue, config.Server.cacheSize)
	server.Monitoring.config = config.Hosts
	server.connection = config.Server.IP + ":" + config.Server.Port

	return server, err
}

func (s *Server) RunServer() {
	l, err := net.Listen("tcp", s.connection)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	agentRequest, _ := getDataFromConn(conn)

	var out []byte

	rsp := ServerResponse{
		Response: "success",
	}

	switch agentRequest.Request {
	case "agent data":
		itemCount := len(agentRequest.Data)
		rsp.Info = fmt.Sprintf("processed: %d; failed: 0; total: %d; seconds spent: 0.000182", itemCount, itemCount)
		for _, item := range agentRequest.Data {
			s.Cache <- item
		}
	case "active checks":
		rsp.Data = s.Monitoring.GetConfig(agentRequest.Host)
	}

	out, _ = json.Marshal(rsp)
	response := generateResponse(out)
	conn.Write(response)
}

func generateResponse(data []byte) (response []byte) {
	rspLen := make([]byte, 4)
	dataLen := uint32(len(data))
	binary.LittleEndian.PutUint32(rspLen, dataLen)
	response = append(HEADER, rspLen...)
	response = append(response, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	response = append(response, data...)
	return
}

func getDataFromConn(conn net.Conn) (agentRequest AgentRequest, err error) {
	header := make([]byte, 5)

	if _, err = conn.Read(header); err != nil {
		return
	}

	protocol_details := int(header[4])
	fieldLen := 4
	if protocol_details&LARGE_PACKET_FLAG != 0 {
		fieldLen = 8
	}

	meta := make([]byte, fieldLen*2)
	conn.Read(meta)

	reqLen := int32(binary.LittleEndian.Uint32(meta[:fieldLen]))
	data := make([]byte, reqLen)

	n, err := conn.Read(data)
	if intRegLen := int(reqLen); n != intRegLen {
		err = fmt.Errorf("got less data (%d) than expected (%d)", n, intRegLen)
	}

	if err != nil {
		return
	}

	json.Unmarshal(data, &agentRequest)
	return
}
