package data

import (
	"encoding/binary"
	"hash/crc32"
)

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// crc type keySize valueSize
// 4 + 1     + 5     + 5 = 15
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

// LogRecord 写入到数据文件的记录
// 之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecord 头部信息
type LogRecordHeader struct {
	crc        uint32        // 校验值
	recordType LogRecordType // 标识类型
	keySize    uint32        // key 长度
	valueSize  uint32        // value 长度
}

// LogRecordPos 数据内存索引，主要是描述数据在磁盘上的位置
type LogRecordPos struct {
	Fid    uint32 // 文件 id, 表示将数据存储到了哪个文件当中
	Offset int64  // 偏移，表示将数据存储到了数据文件中的哪个位置
}

// EncodeLogRecord 对LogRecord进行编码，返回字节数组及长度
func EncodeLogRecord(record *LogRecord) ([]byte, int64) {
	// crc type, keySize, valueSize, key, value
	// 4    1      5         5
	// 初始化一个 header 部分的字节数组

	header := make([]byte, maxLogRecordHeaderSize)

	header[4] = record.Type
	var index = 5

	index += binary.PutVarint(header[index:], int64(len(record.Key)))
	index += binary.PutVarint(header[index:], int64(len(record.Value)))

	var size = index + len(record.Key) + len(record.Value)
	encBytes := make([]byte, size)

	copy(encBytes[:index], header[:index])
	copy(encBytes[index:], record.Key)
	copy(encBytes[index+len(record.Key):], record.Value)

	// 对整个record 的数据进行 crc 校验
	crc := crc32.ChecksumIEEE(encBytes[4:])
	binary.LittleEndian.PutUint32(encBytes[:4], crc)

	return encBytes, int64(size)

}

// 对字节数组中的Header信息进行解码
func decodeLogRecordHeader(buf []byte) (*LogRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &LogRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]),
		recordType: buf[4],
	}

	var index = 5
	// 取出实际的key size
	keySize, n := binary.Varint(buf[index:])
	header.keySize = uint32(keySize)
	index += n

	// 取出实际的 value size
	valueSize, n := binary.Varint(buf[index:])
	header.valueSize = uint32(valueSize)
	index += n

	return header, int64(index)
}

func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	if lr == nil {
		return 0
	}

	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, lr.Key)
	crc = crc32.Update(crc, crc32.IEEETable, lr.Value)

	return crc
}
