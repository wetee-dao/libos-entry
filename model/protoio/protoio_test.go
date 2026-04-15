package protoio

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 定义一个测试用的 proto 消息
type testMessage struct {
	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Data  []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *testMessage) Reset()         { *m = testMessage{} }
func (m *testMessage) String() string { return "" }
func (m *testMessage) ProtoMessage()  {}

// TestMarshalDeterministicConsistency 测试确定性序列化的一致性
func TestMarshalDeterministicConsistency(t *testing.T) {
	msg := &testMessage{
		Name:  "test",
		Value: 12345,
		Data:  []byte("test_data"),
	}

	// 多次序列化应该产生相同的结果
	var hashes [][32]byte
	for i := 0; i < 10; i++ {
		data, err := MarshalDeterministic(msg)
		require.NoError(t, err)
		hash := sha256.Sum256(data)
		hashes = append(hashes, hash)
	}

	// 所有 hash 应该相同
	for i := 1; i < len(hashes); i++ {
		assert.Equal(t, hashes[0], hashes[i], "第 %d 次序列化的 hash 与第一次不同", i)
	}
}

// TestWriteMessageAndReadMsg 测试写入和读取消息
func TestWriteMessageAndReadMsg(t *testing.T) {
	original := &testMessage{
		Name:  "test",
		Value: 12345,
		Data:  []byte("test_data"),
	}

	// 写入消息
	var buf bytes.Buffer
	err := WriteMessage(original, &buf)
	require.NoError(t, err)

	// 读取消息
	reader := NewDelimitedReader(&buf, 1024*1024)
	decoded := &testMessage{}
	_, err = reader.ReadMsg(decoded)
	require.NoError(t, err)

	// 验证内容
	assert.Equal(t, original.Name, decoded.Name)
	assert.Equal(t, original.Value, decoded.Value)
	assert.Equal(t, original.Data, decoded.Data)
}

// TestMarshalDeterministicWithDifferentMessages 测试不同消息产生不同结果
func TestMarshalDeterministicWithDifferentMessages(t *testing.T) {
	msg1 := &testMessage{
		Name:  "msg1",
		Value: 1000,
	}
	msg2 := &testMessage{
		Name:  "msg2",
		Value: 2000,
	}

	data1, err := MarshalDeterministic(msg1)
	require.NoError(t, err)
	data2, err := MarshalDeterministic(msg2)
	require.NoError(t, err)

	hash1 := sha256.Sum256(data1)
	hash2 := sha256.Sum256(data2)

	assert.NotEqual(t, hash1, hash2, "不同消息应该产生不同的 hash")
}

// TestMarshalDeterministicEmptyMessage 测试空消息
func TestMarshalDeterministicEmptyMessage(t *testing.T) {
	msg := &testMessage{}

	data, err := MarshalDeterministic(msg)
	require.NoError(t, err)
	assert.NotNil(t, data)

	// 多次序列化应该一致
	data2, err := MarshalDeterministic(msg)
	require.NoError(t, err)
	assert.Equal(t, data, data2)
}

// TestNewDelimitedWriter 测试 DelimitedWriter
func TestNewDelimitedWriter(t *testing.T) {
	msg := &testMessage{
		Name:  "test",
		Value: 1234,
	}

	var buf bytes.Buffer
	writer := NewDelimitedWriter(&buf)
	defer writer.Close()

	n, err := writer.WriteMsg(msg)
	require.NoError(t, err)
	assert.Greater(t, n, 0)

	// 验证可以读取
	reader := NewDelimitedReader(&buf, 1024*1024)
	decoded := &testMessage{}
	_, err = reader.ReadMsg(decoded)
	require.NoError(t, err)
	assert.Equal(t, msg.Name, decoded.Name)
}

// TestUnmarshalDelimited 测试 UnmarshalDelimited
func TestUnmarshalDelimited(t *testing.T) {
	original := &testMessage{
		Name:  "test",
		Value: 12345,
		Data:  []byte("test_data"),
	}

	// 序列化
	var buf bytes.Buffer
	err := WriteMessage(original, &buf)
	require.NoError(t, err)

	// 使用 UnmarshalDelimited 反序列化
	decoded := &testMessage{}
	err = UnmarshalDelimited(buf.Bytes(), decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Name, decoded.Name)
	assert.Equal(t, original.Value, decoded.Value)
}

// TestReadMessage 测试 ReadMessage
func TestReadMessage(t *testing.T) {
	original := &testMessage{
		Name:  "test",
		Value: 12345,
	}

	// 序列化
	var buf bytes.Buffer
	err := WriteMessage(original, &buf)
	require.NoError(t, err)

	// 使用 ReadMessage 读取
	decoded := &testMessage{}
	err = ReadMessage(&buf, decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Name, decoded.Name)
	assert.Equal(t, original.Value, decoded.Value)
}

// TestMarshalDeterministicWithProtoMessage 测试使用标准 proto.Message
func TestMarshalDeterministicWithProtoMessage(t *testing.T) {
	// 使用一个简单的 proto 消息
	msg := &testMessage{
		Name:  "proto_test",
		Value: 999,
		Data:  []byte("proto_data"),
	}

	data, err := MarshalDeterministic(msg)
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)

	// 验证可以反序列化
	decoded := &testMessage{}
	err = proto.Unmarshal(data, decoded)
	require.NoError(t, err)
	assert.Equal(t, msg.Name, decoded.Name)
	assert.Equal(t, msg.Value, decoded.Value)
	assert.Equal(t, msg.Data, decoded.Data)
}

// BenchmarkMarshalDeterministic 基准测试
func BenchmarkMarshalDeterministic(b *testing.B) {
	msg := &testMessage{
		Name:  "benchmark",
		Value: 1234567890,
		Data:  make([]byte, 1024), // 1KB 数据
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalDeterministic(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestWriteMessageHashConsistency 测试 WriteMessage 的 hash 一致性
func TestWriteMessageHashConsistency(t *testing.T) {
	msg := &testMessage{
		Name:  "test",
		Value: 12345,
		Data:  []byte("test_data"),
	}

	// 多次 WriteMessage 应该产生相同的字节序列
	var hashes [][32]byte
	for i := 0; i < 10; i++ {
		var buf bytes.Buffer
		err := WriteMessage(msg, &buf)
		require.NoError(t, err)

		data := buf.Bytes()
		hash := sha256.Sum256(data)
		hashes = append(hashes, hash)
	}

	// 所有 hash 应该相同
	for i := 1; i < len(hashes); i++ {
		assert.Equal(t, hashes[0], hashes[i], "第 %d 次 WriteMessage 的 hash 与第一次不同", i)
	}
}

// TestWriteMessageAndMarshalDeterministicConsistency 测试 WriteMessage 和 MarshalDeterministic 的一致性
func TestWriteMessageAndMarshalDeterministicConsistency(t *testing.T) {
	msg := &testMessage{
		Name:  "test",
		Value: 12345,
		Data:  []byte("test_data"),
	}

	// 使用 WriteMessage 序列化
	var buf bytes.Buffer
	err := WriteMessage(msg, &buf)
	require.NoError(t, err)
	writeMsgData := buf.Bytes()

	// 使用 MarshalDeterministic 序列化
	marshalData, err := MarshalDeterministic(msg)
	require.NoError(t, err)

	// WriteMessage 的数据应该包含 varint 长度前缀 + MarshalDeterministic 的数据
	// 验证：跳过 varint 前缀后，内容应该与 MarshalDeterministic 一致
	varintLen := 1
	for i := 0; i < len(writeMsgData); i++ {
		if writeMsgData[i] < 0x80 {
			varintLen = i + 1
			break
		}
	}

	// 验证长度前缀正确
	length := uint64(0)
	for i := 0; i < varintLen; i++ {
		length |= uint64(writeMsgData[i]&0x7F) << (7 * i)
	}
	assert.Equal(t, len(marshalData), int(length), "长度前缀应该等于 MarshalDeterministic 的数据长度")

	// 验证内容一致
	assert.Equal(t, marshalData, writeMsgData[varintLen:], "WriteMessage 的内容应该与 MarshalDeterministic 一致")

	// 验证 hash 一致（去掉长度前缀）
	writeMsgHash := sha256.Sum256(writeMsgData[varintLen:])
	marshalHash := sha256.Sum256(marshalData)
	assert.Equal(t, marshalHash, writeMsgHash, "WriteMessage 和 MarshalDeterministic 的 hash 应该一致")
}

// BenchmarkWriteMessage 基准测试
func BenchmarkWriteMessage(b *testing.B) {
	msg := &testMessage{
		Name:  "benchmark",
		Value: 1234567890,
		Data:  make([]byte, 1024),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := WriteMessage(msg, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
