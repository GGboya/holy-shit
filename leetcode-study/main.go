package main

import (
	"encoding/json"
	"fmt"
	"leetcode/ggdb"
)

func main() {
	// lastSubmitTime, err := utils.FetchLastSubmitTime()
	// if err != nil {
	// 	fmt.Println("Error fetching last submit time:", err)
	// 	return
	// }

	// currentTime := time.Now()
	// duration := currentTime.Sub(*lastSubmitTime)
	// fmt.Println("最近提交时间:", lastSubmitTime)
	// fmt.Println("当前时间:", currentTime)
	// fmt.Println("时间差:", duration)

	// if duration < 24*time.Hour {
	// 	fmt.Println("最近一次提交是在24小时内。")
	// } else {
	// 	fmt.Println("最近一次提交超过24小时。")
	// }

	op := ggdb.Options{
		DataFileSize: 1 << 20,
		DirPath:      "./data",
		IndexType:    ggdb.BTree,
		SyncWrites:   true,
	}
	db, err := ggdb.Open(op)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	user2 := User{
		ID:    "selfknow",
		QQ:    "202572197",
		Level: 0,
	}

	value, err := json.Marshal(user2)
	if err != nil {
		fmt.Println("Error marshaling user:", err)
		return
	}
	db.Put([]byte("selfknow"), value)

	for _, key := range db.Iterator() {
		fmt.Println(string(key))
	}

}

type User struct {
	ID    string `json:"id"`
	QQ    string `json:"qq"`
	Level int    `json:"level"`
}
