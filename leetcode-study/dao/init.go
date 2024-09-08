// dao/init.go

package dao

import (
	"leetcode/ggdb"
	"os"
)

var db *ggdb.DB
var dbSecret *ggdb.DB

// Init 初始化数据库
func init() {
	op := ggdb.Options{
		DataFileSize: 1 << 20,
		DirPath:      os.Getenv("DIRPATH"),
		IndexType:    ggdb.BTree,
		SyncWrites:   true,
	}
	var err error
	db, err = ggdb.Open(op)
	if err != nil {
		panic(err)
	}

	dbSecret, err = ggdb.Open(ggdb.Options{
		DataFileSize: 1 << 20,
		DirPath:      "./storge/datasecret",
		IndexType:    ggdb.BTree,
		SyncWrites:   true,
	})

	if err != nil {
		panic(err)
	}
}

// GetDB 获取普通数据库客户端
func NewDB() DBClientInterface {
	return NewDBClient(db)

}

// GetSecretDB 获取密钥数据库客户端
func NewSecretDB() DBClientInterface {
	return NewDBClient(dbSecret)
}
