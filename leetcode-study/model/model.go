package model

type User struct {
	ID     string `json:"id"`
	QQ     string `json:"qq"`
	Level  int    `json:"level"`
	QQName string `json:"qq_name"`
}
