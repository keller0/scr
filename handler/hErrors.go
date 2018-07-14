package handle

var responseErr = map[string]string{

	"Bad Requset":             "4000", // 请求格式错误
	"Email is not valid":      "4001", // 邮箱地址不合法
	"Username is not valid":   "4002", // 用户名不合法
	"Password is too short":   "4003", // 密码太短了
	"Username is too long":    "4004", // 用户名太长了
	"Language not support":    "4005", // 语言不支持
	"Payload not valid":       "4006", // 代码格式错误
	"Too much output":         "4007", // 代码输出太多
	"Time out":                "4080", // 代码超时
	"User Already Exist":      "4090", // 用户已经存在了
	"Email Already Exist":     "4091", // 邮箱地址已经存在了
	"Already Liked":           "4092", // 已经点过赞了
	"Wrong Password":          "4010", // 密码错误
	"Like Code Not Allowed":   "4011", // 对代码点赞需要登录
	"UserNotExist":            "4040", // 用户不存在
	"CodeNotExist":            "4041", // 代码不存在
	"Get Code Not Allowed":    "4030", // 没有权限获取代码
	"Update Code Not Allowed": "4031", // 没有权限更新代码

	"ServerErr Register Failed":    "5001", // 注册失败 服务器错误
	"ServerErr Create Code Failed": "5002", // 创建代码失败 服务器错误
	"ServerErr Get Code Failed":    "5003", // 获取代码失败
	"ServerErr Like Code Failed":   "5004", // 点赞失败
	"Run code error":               "5005", // 运行代码错误
}
