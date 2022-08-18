package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"net"
	"os"
)

var HEADER = []byte{'Z', 'B', 'X', 'D', '\x01'}

type ActiveItem struct {
	Key     string `json:"key"`
	Delay   int32  `json:"delay"`
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

func RunServer(listenIP string, ListenPort string, valueBuffer chan ActiveItemValue) {
	l, err := net.Listen("tcp", listenIP+":"+ListenPort)
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
		go handleRequest(conn, valueBuffer)
	}
}

func handleRequest(conn net.Conn, buffer chan ActiveItemValue) {
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
			buffer <- item
		}
	case "active checks":
		rsp.Data = Monitoring.GetConfig(agentRequest.Host)
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

	fieldLen := 4
	if header[4] != '\x01' {
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
