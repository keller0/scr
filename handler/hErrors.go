package handle

var responseErr = map[string]string{

	"Bad Requset":             "4000", // 请求格式错误
	"Email is not valid":      "4001", // 邮箱地址不合法
	"Username is not valid":   "4002", // 用户名不合法
	"Password is too short":   "4003", // 密码太短了
	"Username is too long":    "4004", // 用户名太长了
	"User Alread Exist":       "4090", // 用户已经存在了
	"Email Alread Exist":      "4091", // 邮箱地址已经存在了
	"Wrong Password":          "4010", // 密码错误
	"UserNotExist":            "4040", // 用户不存在
	"CodeNotExist":            "4001", // 代码不存在
	"Get Code Not Allowed":    "4030", // 没有权限获取代码
	"Update Code Not Allowed": "4031", // 没有权限更新代码

	"ServerErr Register Failed":    "5001", // 注册失败 服务器错误
	"ServerErr Create Code Failed": "5002", // 创建代码失败 服务器错误
	"ServerErr Get Code Failed":    "5003", // 获取代码失败
}
