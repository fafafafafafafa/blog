package g

import "fmt"

const (
	SUCCESS = 0   // 成功业务码
	FAIL    = 500 // 失败业务码
)

// 自定义业务 error 类型
type Result struct {
	code int
	msg  string
}

func (e Result) Code() int {
	return e.code
}

func (e Result) Msg() string {
	return e.msg
}

var (
	_codes    = map[int]struct{}{}   // 注册过的错误码集合, 防止重复
	_messages = make(map[int]string) // 根据错误码获取错误信息
)

// 注册一个响应码, 不允许重复注册
func RegisterResult(code int, msg string) Result {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	if msg == "" {
		panic("错误码消息不能为空")
	}

	_codes[code] = struct{}{}
	_messages[code] = msg

	return Result{
		code: code,
		msg:  msg,
	}
}

// 根据响应码获取响应信息
func GetMsg(code int) string {
	return _messages[code]
}

var (
	OkResult   = RegisterResult(SUCCESS, "OK")
	FailResult = RegisterResult(FAIL, "FAIL")
)

var (
	ErrRequest = RegisterResult(9001, "请求参数格式错误")
	ErrDbOp    = RegisterResult(9004, "数据库操作异常")

	ErrPassword     = RegisterResult(1002, "密码错误")
	ErrUserNotExist = RegisterResult(1003, "该用户不存在")

	ErrTokenCreate = RegisterResult(1205, "TOKEN 生成失败")
)
