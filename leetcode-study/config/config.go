package config

import "os"

const (
	Authority = "leetcode.cn"

	UserAgent   = "apifox/1.0.0 (https://www.apifox.cn)"
	ContentType = "application/json"
)

var (
	Cookie     = os.Getenv("COOKIE")
	XCSRFToken = os.Getenv("XCSRFTOKEN")

	Headers = map[string]string{
		"authority":     Authority,
		"authorization": "",
		"cookie":        Cookie,
		"x-csrftoken":   XCSRFToken,
		"User-Agent":    UserAgent,
		"content-type":  ContentType,
	}
	Lazy = map[int]string{
		0: "平民",
		1: "懒惰者",
		2: "懒惰师",
		3: "懒惰大师",
		4: "懒惰王",
		5: "懒惰皇",
		6: "懒惰半圣",
		7: "懒惰圣人",
		8: "懒惰帝",
		9: "飞升",
	}
)
