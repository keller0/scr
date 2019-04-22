package handler

var responseErr = map[string]string{

	"Bad Requset":           "4000", // 请求格式错误
	"Password is too short": "4003", // 密码太短了
	"Language not support":  "4005", // 语言不支持
	"Payload not valid":     "4006", // 代码格式错误
	"Too much output":       "4007", // 代码输出太多
	"Time out":              "4080", // 代码超时

	"Run code error": "5005", // 运行代码错误
}
