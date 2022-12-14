package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"GoUserManaSys/log"
	"GoUserManaSys/rpc"
	"GoUserManaSys/utils"
)

var (
	//rpc _client
	_client rpc.Client
	//template parameters
	_loginT   *template.Template
	_profileT *template.Template
	_jumpT    *template.Template
	_addT     *template.Template
)

//parse html
func init() {
	_loginT = template.Must(template.ParseFiles("./web/login.html"))
	_profileT = template.Must(template.ParseFiles("./web/profile.html"))
	_jumpT = template.Must(template.ParseFiles("./web/jump.html"))
	_addT = template.Must(template.ParseFiles("./web/add.html"))
}

func main() {
	//init log
	if err := log.ConfigLog(utils.HTTPServerLogPath, log.LevelInfo); err != nil {
		panic(err)
	}
	var err error
	//creat rpc _client connect pool
	_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
	if err != nil {
		fmt.Println(err)
	}
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(utils.StaticFilePath))))

	//assemble handle func for http requests
	http.HandleFunc("/AddUser", AddUser)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/GetInfo", GetInfo)
	http.HandleFunc("/UpdateNickName", UpdateNickName)
	http.HandleFunc("/UploadPic", UploadPic)
	http.HandleFunc("/Logout", Logout)
	//turn on http listen and serve
	//http.ListenAndServe(utils.HTTPServerPort, nil)
	go http.ListenAndServe(utils.HTTPServerPort, nil)
	fmt.Println("http server started", utils.HTTPServerPort)
	gracefulExit()
}

func gracefulExit() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch
	log.InfoLog("got a signal" + sig.String())
	fmt.Println("handling shutdown")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_client.Close()
	log.InfoLog("[INFO] ------exited--------")
	fmt.Println("shutdown successfully")
}

//handle the login request
func Login(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		userName := req.FormValue("username")
		userName = template.HTMLEscapeString(userName)
		passWord := req.FormValue("password")
		passWord = template.HTMLEscapeString(passWord)
		if userName == "" || passWord == "" {
			templateLogin(res, utils.MsgLogin{Msg: "??????????????????????????????"})
			return
		}
		//request
		req := utils.ReqLogin{
			UserName: userName,
			PassWord: passWord,
		}
		//response
		rsp := utils.ResLogin{}
		err := _client.Call("Login", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("Login", req, &rsp); err != nil {
				log.ErrorLog("http_server_login: call failed.username:%s,err:%s", userName, err)
			}

		}
		if rsp.Code == utils.Success {
			//if login success,then send token as cookie to http
			cookie := http.Cookie{Name: "token", Value: rsp.Token, MaxAge: utils.TokenLife}
			http.SetCookie(res, &cookie)
			//jump html and send message
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "???????????????"})
		} else {
			//username wrong, password wrong ,login failed, et.al
			templateLogin(res, utils.MsgLogin{Msg: utils.GetErrMsg(rsp.Code)})
		}
		//log.InfoLog("http_server_login: username:%s login code: %d", userName, rsp.Code)
	}
	return
}

//handle logout request
func Logout(res http.ResponseWriter, req *http.Request) {
	//??????token
	if req.Method == "POST" {
		token, err := req.Cookie("token")
		if err != nil {
			templateLogin(res, utils.MsgLogin{Msg: ""})
		}
		userName := req.FormValue("username")
		userName = template.HTMLEscapeString(userName)
		req := utils.ReqLogout{
			UserName: userName,
			Token:    token.Value,
		}
		rsp := utils.ResLogout{}
		err = _client.Call("Logout", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("Logout", req, &rsp); err != nil {
				log.ErrorLog("http_server_logout: call failed.username:%s,err:%s", userName, err)
			}
		}
		templateLogin(res, utils.MsgLogin{Msg: "????????????"})
		//log.InfoLog("http_server_logout: logout.username:%s,err:%s", userName, err)
	}
}

//handle getinfo request
func GetInfo(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		//get token from cookie
		//token, err := req.Cookie("token")
		var err error
		if err != nil {
			templateLogin(res, utils.MsgLogin{Msg: ""})
			return
		}
		userName := req.FormValue("username")
		userName = template.HTMLEscapeString(userName)
		req := utils.ReqGetInfo{
			UserName: userName,
			//Token:    token.Value,
			Token: utils.TestToken,
		}
		rsp := utils.ResGetInfo{}
		err = _client.Call("GetInfo", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("GetInfo", req, &rsp); err != nil {
				log.ErrorLog("http_server_GetInfo: call failed.username:%s,err:%s", userName, err)
			}

		}
		//getting success
		if rsp.Code == utils.Success {
			//if no image, then use the default image
			if rsp.ProfilePicture == "" {
				rsp.ProfilePicture = utils.DefaultImage
			}
			//log.InfoLog("http_server_getIfo: getInfo success username:%s", userName)
			//display information in html
			templateProfile(res, utils.MsgGetInfo{
				UserName:       rsp.UserName,
				NickName:       rsp.NickName,
				ProfilePicture: rsp.ProfilePicture})
			return
		}
		//errors
		switch rsp.Code {
		case utils.ErrUserNotExit:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      utils.GetErrMsg(rsp.Code)})
		case utils.Err:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      utils.GetErrMsg(rsp.Code)})
		case utils.ErrTokenWrong:
			templateLogin(res, utils.MsgLogin{Msg: "??????????????????"})
		}
		//log.ErrorLog("http_server_getIfo: getInfo failed. username:%s getInfo code: %d", userName, rsp.Code)
	}
	return
}

//handle update nickname request
func UpdateNickName(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		//get token from cookie
		//token, err := req.Cookie("token")
		var err error
		if err != nil {
			fmt.Println("UpdateNickname,get token failed")
			templateLogin(res, utils.MsgLogin{Msg: ""})
			return
		}
		userName := req.FormValue("username")
		userName = template.HTMLEscapeString(userName)
		nickName := req.FormValue("nickname")
		nickName = template.HTMLEscapeString(nickName)
		req := utils.ReqUpdNickName{
			UserName: userName,
			NickName: nickName,
			//Token:    token.Value,
			Token: utils.TestToken,
		}
		rsp := utils.ResUpdNickName{}
		err = _client.Call("UpdateNickName", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("UpdateNickName", req, &rsp); err != nil {
				log.ErrorLog("http_server_UpdateNickName: call failed.username:%s,err:%s", userName, err)
			}
		}
		switch rsp.Code {
		case utils.ErrUserNotExit:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "???????????????"})
		case utils.Err:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
		case utils.Success:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "????????????"})
		case utils.ErrTokenWrong:
			templateLogin(res, utils.MsgLogin{Msg: "????????????"})
		default:
			templateLogin(res, utils.MsgLogin{Msg: "????????????"})
		}
		//log.InfoLog("http_server_UpdateNickName: username:%s login code: %d", userName, rsp.Code)

	}

}

//handle upload profile picture request
func UploadPic(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		// get token from cookie
		//token, err := req.Cookie("token")
		var err error
		if err != nil {
			fmt.Println("UploadPic,get token failed")
			templateLogin(res, utils.MsgLogin{Msg: ""})
			return
		}
		userName := req.FormValue("username")
		userName = template.HTMLEscapeString(userName)
		//get image from html
		//file, head, err := req.FormFile("image")
		file, err := os.Open(utils.StaticFilePath + utils.DefaultImage)
		//loadImage ,err :=
		if err != nil {
			fmt.Println(err)
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
			return
		}
		//check whether the file format is jpg,jpeg,png,gif
		if utils.CheckImage(file.Name()) == false {
			templateJump(res, utils.MsgJump{Msg: "???????????????????????????"})
			return
		}
		//rename the image to avoid name repeat
		newName := utils.NewImgName(file.Name())
		//test image save address
		filepath := "./statictest/" + newName
		showFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
		defer showFile.Close()
		//add the image to the filepath
		_, err = io.Copy(showFile, file)
		if err != nil {
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
		}
		//
		req := utils.ReqUploadPic{
			UserName: userName,
			Picture:  newName,
			//Token:    token.Value,
			Token: utils.TestToken,
		}
		rsp := utils.ResUploadPic{}
		err = _client.Call("UploadPic", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("UploadPic", req, &rsp); err != nil {
				log.ErrorLog("http_server_UploadPic: call failed.username:%s,err:%s", userName, err)
			}
		}
		switch rsp.Code {
		case utils.Success:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
		case utils.Err:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
		case utils.ErrTokenWrong:
			templateLogin(res, utils.MsgLogin{
				Msg: "???????????????"})
		case utils.ErrUserNotExit:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "???????????????"})
		default:
			templateJump(res, utils.MsgJump{
				UserName: userName,
				Msg:      "??????????????????"})
		}
		//log.InfoLog("http_server_UploadPic: username:%s login code: %d", userName, rsp.Code)
	}
	return
}

//handle add user requst
func AddUser(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		userName := req.FormValue("username")
		fmt.Println(userName)
		//use template.HTMLEscapeString to handle special symbols '"<>_ to avoid sql injection
		userName = template.HTMLEscapeString(userName)
		fmt.Println(userName)
		passWord := req.FormValue("password")
		passWord = template.HTMLEscapeString(passWord)
		nickName := req.FormValue("nickname")
		nickName = template.HTMLEscapeString(nickName)
		if userName == "" || passWord == "" {
			templateAdd(res, utils.MsgAdd{Msg: "??????????????????????????????"})
			return
		}
		req := utils.ReqAdd{
			UserName: userName,
			PassWord: passWord,
			NickName: nickName,
		}
		rsp := utils.ResAdd{}
		err := _client.Call("AddUser", req, &rsp)
		if err != nil {
			fmt.Println("Disconnected, re-establishing connection...")
			//close pool
			_client.Close()
			//re-creat rpc client
			_client, err = rpc.NewClient(utils.ClientPoolSize, utils.ServerPort)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully re-established connection...")
			if err = _client.Call("AddUser", req, &rsp); err != nil {
				log.ErrorLog("http_server_add: call failed.username:%s,err:%s", userName, err)
			}
		}
		//display front-end page and jump
		switch rsp.Code {
		case utils.Success:
			templateLogin(res, utils.MsgLogin{Msg: "???????????????????????????"})
		case utils.ErrUserExit:
			templateAdd(res, utils.MsgAdd{Msg: "?????????????????????"})
		case utils.ErrNil:
			templateAdd(res, utils.MsgAdd{Msg: "??????????????????????????????"})
		}
		log.InfoLog("http_server_add: username:%s add code: %d", userName, rsp.Code)
	}

}

//http login page.
func templateLogin(rw http.ResponseWriter, resp utils.MsgLogin) {
	if err := _loginT.Execute(rw, resp); err != nil {
		fmt.Println(err)
	}
}

//http jump page.
func templateJump(rw http.ResponseWriter, rsp utils.MsgJump) {
	if err := _jumpT.Execute(rw, rsp); err != nil {
		fmt.Println(err)
	}
}

//http profile page.
func templateProfile(rw http.ResponseWriter, resp utils.MsgGetInfo) {
	if err := _profileT.Execute(rw, resp); err != nil {
		fmt.Println(err)
	}
}

//http add page.
func templateAdd(rw http.ResponseWriter, resp utils.MsgAdd) {
	if err := _addT.Execute(rw, resp); err != nil {
		fmt.Println(err)
	}
}
