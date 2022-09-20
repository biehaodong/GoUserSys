package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"GoUserManaSys/dao"
	"GoUserManaSys/log"
	"GoUserManaSys/rpc"
	"GoUserManaSys/utils"
)

func main() {
	//init log.
	if err := log.ConfigLog(utils.TCPServerLogPath, log.LevelInfo); err != nil {
		panic(err)
	}
	//init rpc server
	s := rpc.NewServer()
	//register server
	s.Register("AddUser", AddUser, AddUserServ)
	s.Register("Login", Login, LoginServ)
	s.Register("GetInfo", GetInfo, GetInfoServ)
	s.Register("UpdateNickName", UpdateNickName, UpdateNickNameServ)
	s.Register("UploadPic", UploadPic, UploadPicServ)
	s.Register("Logout", Logout, LogoutServ)
	//listen
	l, err := s.Listen(utils.ServerPort)
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("handling shutdown")
		s.Shutdown(l)
		dao.CloseRedis()
		fmt.Println("shutdown successfully")
	}()
	if err != nil {
		fmt.Println(err)
	}
	//serve
	s.Serve(l)
}

//add user interface
func AddUser(f interface{}) interface{} {
	return AddUserServ(*f.(*utils.ReqAdd))
}

// Login interface.
func Login(f interface{}) interface{} {
	return LoginServ(*f.(*utils.ReqLogin))
}

//logout interface
func Logout(f interface{}) interface{} {
	return LogoutServ(*f.(*utils.ReqLogout))
}

//GetInfo
func GetInfo(f interface{}) interface{} {
	return GetInfoServ(*f.(*utils.ReqGetInfo))
}

//update profile picture interface
func UploadPic(f interface{}) interface{} {
	return UploadPicServ(*f.(*utils.ReqUploadPic))
}

//update nicknam interface
func UpdateNickName(f interface{}) interface{} {
	return UpdateNickNameServ(*f.(*utils.ReqUpdNickName))
}

//login service
func LoginServ(req utils.ReqLogin) (res utils.ResLogin) {

	code := dao.Login(req.UserName, req.PassWord)
	res.Code = code

	//not success, then return
	if code != utils.Success {
		log.ErrorLog("tcp_server_login: login failed. username:%s, code:%d", req.UserName, code)
		return
	}
	//+++++++++++++++++++++++++++++++++++change release token or test token+++++++++++++++++
	//release token
	token := utils.GetToken()
	//test token
	//token := utils.GetTokenTest()
	//+++++++++++++++++++++++++++++++++++change release token or test token+++++++++++++++++
	//set name->token in redis
	rCode := dao.SetToken(req.UserName, token, int64(utils.TokenLife))
	if rCode != utils.Success {
		//set token to redis wrong, return
		log.ErrorLog("tcp_server_login: redis set failed. username:%s", req.UserName)
		res.Code = utils.ErrRedisSet
		return
	}

	//Set success
	res.Token = token
	res.UserName = req.UserName
	//log.InfoLog("tcp_server_login: login success. username:%s", req.UserName)
	return
}

//add user service
func AddUserServ(req utils.ReqAdd) (res utils.ResAdd) {
	//username or password can't be nil
	if req.UserName == "" || req.PassWord == "" {
		res.Code = utils.ErrNil
		return
	}
	code := dao.AddUser(req.UserName, req.PassWord, req.NickName)
	res.Code = code

	if code != utils.Success {
		log.ErrorLog("tcp_server_add: add failed. username:%s, code:%d", req.UserName, code)
		return
	}
	log.InfoLog("tcp_server_add: add success. username:%s", req.UserName)
	return
}

//get info service
func GetInfoServ(req utils.ReqGetInfo) (res utils.ResGetInfo) {
	//if token don't match ,return
	code := dao.CheckToken(req.UserName, req.Token)
	if code != utils.Success {
		log.ErrorLog("tcp_server_getInfo: getInfo failed. username:%s, code:%d", req.UserName, code)
		res.Code = code
		return
	}
	//token right, get data from redis first
	nickname, profilepicture, hasData, errCode := dao.RdsGetInfo(req.UserName)
	//err
	if errCode == utils.ErrRedisGet {
		log.ErrorLog("tcp_server_getInfo: getInfo from redis failed. username:%s", req.UserName)
		res.Code = utils.ErrRedisGet
		//return
	}
	//if data exist
	if hasData {
		res.Code = utils.Success
		res.NickName = nickname
		res.UserName = req.UserName
		res.ProfilePicture = profilepicture
		//fmt.Println("从redis获取信息成功")
		//fmt.Println("tcp_server_getInfo: getInfo success. username:%s", req.UserName)
		//log.InfoLog("tcp_server_getInfo: getInfo success. username:%s", req.UserName)
		return
	}
	//redis no data, get data from db
	nickname, profilepicture, errCode = dao.GetInfo(req.UserName)
	if errCode != utils.Success {
		res.Code = code
		log.ErrorLog("tcp_server_getInfo: getInfo from db failed. username:%s", req.UserName)
		return
	}
	//get success from db,set data to redis
	dao.RdsSetInfo(req.UserName, nickname, profilepicture, int64(utils.KeyLife))
	res.Code = utils.Success
	res.NickName = nickname
	res.UserName = req.UserName
	res.ProfilePicture = profilepicture
	//fmt.Println("tcp_server_getInfo: set info to redis success. username:%s", req.UserName)
	//log.InfoLog("tcp_server_getInfo: set info to redis success. username:%s", req.UserName)
	return

}

//update nickname service
func UpdateNickNameServ(req utils.ReqUpdNickName) (res utils.ResUpdNickName) {
	code := dao.CheckToken(req.UserName, req.Token)
	if code != utils.Success {
		res.Code = code
		log.ErrorLog("tcp_server_updateNickname: update failed. username:%s,code:%d", req.UserName, code)
		return
	}
	//token success, invalid redis data
	//code = dao.Invalid(req.UserName)
	//if code != utils.Success {
	//	res.Code = utils.ErrRedisSet
	//	log.ErrorLog("tcp_server_updateNickname: redis set invalid failed. username:%s", req.UserName)
	//	return
	//}
	dao.RdsDelInfo(req.UserName)
	//update db
	code = dao.UpdateNickName(req.UserName, req.NickName)
	res.Code = code
	//fmt.Println("tcp_server_updateNickname: update success. username:%s", req.UserName)
	//log.InfoLog("tcp_server_updateNickname: update success. username:%s", req.UserName)
	return
}

//upload profile picture
func UploadPicServ(req utils.ReqUploadPic) (res utils.ResUploadPic) {
	code := dao.CheckToken(req.UserName, req.Token)
	if code != utils.Success {
		res.Code = code
		log.ErrorLog("tcp_server_uploadPic: uploadPic failed. username:%s,code:%d", req.UserName, code)
		return
	}
	//token success, invalid redis data
	code = dao.Invalid(req.UserName)
	if code != utils.Success {
		//if err, return
		res.Code = utils.ErrRedisSet
		log.ErrorLog("tcp_server_uploadPic: redis set valid failed. username:%s", req.UserName)
		return
	}
	//update db
	code = dao.UploadPic(req.UserName, req.Picture)
	res.Code = code
	fmt.Println("tcp_server_uploadPic: upload success. username:", req.UserName)
	//log.InfoLog("tcp_server_uploadPic: upload success. username:%s", req.UserName)
	return
}

//logout
func LogoutServ(req utils.ReqLogout) (res utils.ResLogout) {
	code := dao.CheckToken(req.UserName, req.Token)
	if code != utils.Success {
		res.Code = code
		log.ErrorLog("tcp_server_Logout: logout failed. username:%s,code:%d", req.UserName, code)
		return
	}
	//token success, invalid redis data
	rCode := dao.SetToken(req.UserName, "", 0)
	if rCode != utils.Success {
		log.ErrorLog("tcp_server_logout: redis set failed. username:%s", req.UserName)
		return
	}
	log.InfoLog("tcp_server_logout: logout success. username:%s", req.UserName)
	return
}
