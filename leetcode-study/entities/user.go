package entities

// 定义数据结构
type User struct {
	ID       string `json:"id"`
	QQ       string `json:"qq"`
	Level    int
	QQName   string `json:"qq_name"`
	NickName string `json:"nickname"`
}

type UserSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
