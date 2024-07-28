package fio

const DataFilePerm = 0644

type IOManager interface {
	// Read 从文件给定的位置读取对应的数据
	Read([]byte, int64) (int, error)

	// Write 写入字节数组到文件中
	Write([]byte) (int, error)

	// Sync 持久化数据
	Sync() error

	// Close 关闭文件
	Close() error

	// Size 获取文件大小
	Size() (int64, error)
}

// NewIOManager 初始化IOManager, 目前支持标准FileID.这里的fileName是文件路径+后缀名
func NewIOManager(fileName string) (IOManager, error) {
	return NewFileIOManager(fileName)
}
