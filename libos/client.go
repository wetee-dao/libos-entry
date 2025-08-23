package libos

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/libos-entry/model"
	"github.com/wetee-dao/libos-entry/model/protoio"
)

func NewTEEClient(addr string) (*TEEClinet, []byte, error) {
	tc := TEEClinet{}
	conn, err := net.DialTimeout("tcp", addr, 20*time.Second)
	if err != nil {
		return nil, nil, err
	}

	_, msg, err := ReadFromApi(conn)
	if err != nil {
		return nil, nil, err
	}

	val := new(model.ApiResp)
	err = protoio.ReadMessage(bytes.NewBuffer(msg), val)
	if err != nil {
		return nil, nil, err
	}

	tc.conn = conn
	tc.requests = make(map[uint64]chan *model.ApiResp)
	return &tc, val.Data, nil
}

// Worker 请求通道
type TEEClinet struct {
	conn     net.Conn
	mu       sync.Mutex
	requests map[uint64]chan *model.ApiResp
}

// Start 启动客户端
func (w *TEEClinet) Start() error {
	conn := w.conn

	for {
		id, msg, err := ReadFromApi(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Disconnected from server:", err)
				return err
			}
			continue
		}

		val := new(model.ApiResp)
		err = protoio.ReadMessage(bytes.NewBuffer(msg), val)
		if err != nil {
			fmt.Println("Client ReadMessage failed:", err)
			continue
		}

		respChan := w.requests[id]
		if respChan != nil {
			respChan <- val
		}
	}
}

// Invoke worker function
func (w *TEEClinet) Invoke(url string, data []byte) ([]byte, error) {
	req := &model.ApiReq{
		Url:  []byte(url),
		Data: data,
	}
	id := GenerateUniqueID()

	buf := new(bytes.Buffer)
	err := types.WriteMessage(req, buf)
	if err != nil {
		return nil, err
	}

	bt, err := Encode(id, buf.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = w.conn.Write(bt)
	if err != nil {
		return nil, err
	}

	// 创建返回 channel
	respChan := make(chan *model.ApiResp, 1)
	w.mu.Lock()
	w.requests[id] = respChan
	w.mu.Unlock()

	resp := <-w.requests[id]
	w.mu.Lock()
	delete(w.requests, id)
	w.mu.Unlock()

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, msg: %s", resp.Code, string(resp.Data))
	}

	return resp.Data, nil
}

// Close 关闭连接
func (w *TEEClinet) Close() error {
	return w.conn.Close()
}

// GenerateUniqueID 生成唯一ID
func GenerateUniqueID() uint64 {
	now := uint64(time.Now().Unix())
	var randPart uint32
	binary.Read(rand.Reader, binary.LittleEndian, &randPart)

	return (now << 32) | uint64(randPart)
}
