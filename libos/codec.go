package libos

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/panjf2000/gnet/v2"
)

const (
	headSize = 8
	bodySize = 4
)

// Codec Protocol format:
//
// * 0           8                       12
// * +-----------+-----------------------+
// * |   head    |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+
func Encode(id uint64, buf []byte) ([]byte, error) {
	bodyOffset := headSize + bodySize
	msgLen := bodyOffset + len(buf)

	data := make([]byte, msgLen)
	idBytes := make([]byte, headSize)
	binary.BigEndian.PutUint64(idBytes, id)

	copy(data, idBytes)

	binary.BigEndian.PutUint32(data[headSize:bodyOffset], uint32(len(buf)))
	copy(data[bodyOffset:msgLen], buf)
	return data, nil
}

func Decode(c gnet.Conn) (uint64, []byte, error) {
	bodyOffset := headSize + bodySize
	buf, err := c.Peek(bodyOffset)
	if err != nil {
		return 0, nil, err
	}

	id := binary.BigEndian.Uint64(buf[:headSize])

	bodyLen := binary.BigEndian.Uint32(buf[headSize:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	buf, err = c.Peek(msgLen)
	if err != nil {
		return 0, nil, err
	}
	body := make([]byte, bodyLen)
	copy(body, buf[bodyOffset:msgLen])
	_, _ = c.Discard(msgLen)

	return id, body, nil
}

func ReadFromApi(c net.Conn) (uint64, []byte, error) {
	bodyOffset := headSize + bodySize
	headerData := make([]byte, bodyOffset)
	_, readTagError := io.ReadFull(c, headerData)
	if readTagError != nil {
		return 0, nil, readTagError
	}

	id := binary.BigEndian.Uint64(headerData[:headSize])

	bodyLen := binary.BigEndian.Uint32(headerData[headSize:bodyOffset])
	bodyData := make([]byte, bodyLen)
	_, readTagError = io.ReadFull(c, bodyData)
	if readTagError != nil {
		return id, nil, readTagError
	}

	return id, bodyData, nil
}
