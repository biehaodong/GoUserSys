package utils

//err category
const (
	Success = 200
	Err     = 500
	//user error information
	ErrUserExit     = 1001
	ErrUserNotExit  = 1002
	ErrPwdWrong     = 1003
	ErrTokenExit    = 1004
	ErrTokenRuntime = 1005
	ErrTokenWrong   = 1006
	ErrNil          = 1007
	ErrRedisSet     = 1008
	ErrRedisGet     = 1009
	//ErrRdsSet       = 1009
)

//err information map
var codeMsg = map[int]string{
	Success:         "OK",
	Err:             "FAIL",
	ErrUserExit:     "用户已存在！",
	ErrPwdWrong:     "密码错误",
	ErrUserNotExit:  "用户不存在",
	ErrTokenExit:    "Token存在",
	ErrTokenRuntime: "token已过期",
	ErrTokenWrong:   "token不正确",
	ErrNil:          "用户名或密码为空",
	ErrRedisSet:     "添加token错误",
	ErrRedisGet:     "获取token错误",
}

//return err information
func GetErrMsg(code int) string {
	return codeMsg[code]
}
