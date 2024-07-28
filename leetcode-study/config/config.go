package config

const (
	Authority   = "leetcode.cn"
	Cookie      = "_bl_uid=32lgt5Iq4estR4ovpv48fej2UI58; gr_user_id=9cd64ed3-db69-4693-b361-a68ed8316706; csrftoken=PPT5fcdFRpPnRbcHzNYbBhwpTv1HxSKtrqyurg4uKxuW0aBUPMTLb0Cu5LwyCmm6; a2873925c34ecbd2_gr_last_sent_cs1=sanxiconze-2; _ga=GA1.2.827755954.1656827815; aliyungf_tc=3f88358c46c63b3cf747bfa7f9c952beb3e2cc9b91a526c574c455d381d3d562; Hm_lpvt_fa218a3ff7179639febdb15e372f411c=1656848229; Hm_lvt_fa218a3ff7179639febdb15e372f411c=1656827815,1656848229; a2873925c34ecbd2_gr_session_id=0f3c7ce7-f057-4e09-9e0d-d8b006ba0375; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=0f3c7ce7-f057-4e09-9e0d-d8b006ba0375; a2873925c34ecbd2_gr_session_id_0f3c7ce7-f057-4e09-9e0d-d8b006ba0375=true; a2873925c34ecbd2_gr_cs1=sanxiconze-2; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMzk3NjA5IiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiYWxsYXV0aC5hY2NvdW50LmF1dGhfYmFja2VuZHMuQXV0aGVudGljYXRpb25CYWNrZW5kIiwiX2F1dGhfdXNlcl9oYXNoIjoiN2RiZjJiMjI5ZDZiMmQ0MWI0ZTMzOTVlYmExMTE5YzI2ZGI0M2Q1ZTE1M2FkNjI1YmJmNDVkYzEwNWI1OWFjMSIsImlkIjozOTc2MDksImVtYWlsIjoic2FueGljb256ZUBnbWFpbC5jb20iLCJ1c2VybmFtZSI6InNhbnhpY29uemUtMiIsInVzZXJfc2x1ZyI6InNhbnhpY29uemUtMiIsImF2YXRhciI6Imh0dHBzOi8vYXNzZXRzLmxlZXRjb2RlLmNuL2FsaXl1bi1sYy11cGxvYWQvZGVmYXVsdF9hdmF0YXIucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2NTY4MjMxNTkuOTk4MTUzNywiZXhwaXJlZF90aW1lXyI6MTY1OTM4MDQwMCwidmVyc2lvbl9rZXlfIjowLCJsYXRlc3RfdGltZXN0YW1wXyI6MTY1Njk0MDY3MX0.-8yAB0-g6Vt2xFGdvuO5HxkV9TSHAywGzeD-Ea2fhE8"
	XCSRFToken  = "PPT5fcdFRpPnRbcHzNYbBhwpTv1HxSKtrqyurg4uKxuW0aBUPMTLb0Cu5LwyCmm6"
	UserAgent   = "apifox/1.0.0 (https://www.apifox.cn)"
	ContentType = "application/json"
)

var (
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
