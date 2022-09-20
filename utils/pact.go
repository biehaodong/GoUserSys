package utils

//add request
type ReqAdd struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	NickName string `json:"nickname"`
}

//add response
type ResAdd struct {
	Code int `json:"code"`
}

//send msg to add.html
type MsgAdd struct {
	Msg string
}

//login request
type ReqLogin struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

//login response
type ResLogin struct {
	UserName string `json:"username"`
	Code     int    `json:"code"`
	Token    string `json:"token"`
}

// send msg to login.html
type MsgLogin struct {
	Msg string
}

//getinfo request
type ReqGetInfo struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
}

//getinfo response
type ResGetInfo struct {
	Code           int    `json:"code"`
	UserName       string `json:"username"`
	NickName       string `json:"nickname"`
	ProfilePicture string `json:"profilepicture"`
}

//send msg
type MsgGetInfo struct {
	UserName       string
	NickName       string
	ProfilePicture string
}

//update nickname requset
type ReqUpdNickName struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Token    string `json:"token"`
}

//update nickname response
type ResUpdNickName struct {
	Code int `json:"code"`
}

//send msg
type MsgUpdNickName struct {
	Msg string
}

//upload picture request
type ReqUploadPic struct {
	UserName string `json:"username"`
	Picture  string `json:"picture"`
	Token    string `json:"token"`
}

//upload picture response
type ResUploadPic struct {
	Code int `json:"code"`
}

//send msg
type MsgUploadPic struct {
	Msg string
}

// sedn msg to jump.html
type MsgJump struct {
	UserName string `json:"username"`
	Msg      string
}

//logout request
type ReqLogout struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
}

//logour response
type ResLogout struct {
	Token string `json:"token"`
	Code  int    `json:"code"`
}

//send msg
type MsgLogout struct {
	Msg string
}
