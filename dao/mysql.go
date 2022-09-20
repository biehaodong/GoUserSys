package dao

import (
	"database/sql"
	"fmt"

	"GoUserManaSys/utils"

	_ "github.com/go-sql-driver/mysql"
)

//use args sql to avoid sql injection
var (
	_addUser        *sql.Stmt
	_updateNickName *sql.Stmt
	_getInfo        *sql.Stmt
	_uploadPic      *sql.Stmt
	_findUser       *sql.Stmt
	_updatePwd      *sql.Stmt
)

//connect to mysql
func init() {
	db, err := sql.Open(utils.Db, utils.DbAddress)
	if err != nil {
		fmt.Println("init mysql failed", err)
	}
	//max amount of the idle connection in pool
	db.SetMaxIdleConns(utils.MaxIdleConns)
	//max amount of the open connections of db
	db.SetMaxOpenConns(utils.MaxOpenConns)
	//max reusable time of connection
	db.SetConnMaxLifetime(utils.ConnMaxLifetime)
	if err = db.Ping(); err != nil {
		fmt.Println("init mysql failed")
	}
	//prepare sql
	//add user
	_addUser = sqlPrepare(db, "INSERT INTO users (user_name, pass_word, nick_name) values (?, ?, ?)")
	//_addUser, _ = db.Prepare("INSERT INTO users (user_name, pass_word) values (?, ?)")
	//update nickName
	_updateNickName = sqlPrepare(db, "UPDATE users SET nick_name = ? WHERE user_name = ?")
	//get nickName and profile picture
	_getInfo = sqlPrepare(db, "SELECT nick_name, profile_picture FROM users WHERE user_name = ?")
	//update profile picture
	_uploadPic = sqlPrepare(db, "UPDATE users SET profile_picture = ? WHERE user_name = ?")
	//check whether user exist
	_findUser = sqlPrepare(db, "SELECT pass_word FROM users WHERE user_name = ?")
	//update password
	_updatePwd = sqlPrepare(db, "UPDATE users SET pass_word = ? WHERE user_name = ?")
	fmt.Println("init mysql success...")
}

//handle prepare errors
func sqlPrepare(db *sql.DB, n string) *sql.Stmt {
	//var code int
	stmt, err := db.Prepare(n)
	if err != nil {
		panic(err)
	}
	return stmt
}

//check whether user exist
func FindUser(name string) int {
	//fmt.Println("find user...")
	rows, err := _findUser.Query(name)
	if err != nil { //error
		fmt.Println(err.Error())
		return utils.Err
	}
	defer rows.Close()

	for rows.Next() { //用户存在
		return utils.ErrUserExit
	}
	//not exist
	return utils.ErrUserNotExit
}

//add user
func AddUser(name string, pwd string, nickname string) int {
	code := FindUser(name)
	if code == utils.Err {
		return utils.Err
	}
	if code == utils.ErrUserExit {
		return utils.ErrUserExit
	}
	//use md5 algorithm to encrypt password
	p := utils.Md5(pwd)
	_, err := _addUser.Exec(name, p, nickname)
	if err != nil {
		return utils.Err
	}
	return utils.Success
}

//update nickName
func UpdateNickName(name string, nickname string) int {
	code := FindUser(name)
	if code == utils.Err {
		return utils.Err
	}
	if code == utils.ErrUserNotExit {
		return utils.ErrUserNotExit
	}
	_, err := _updateNickName.Exec(nickname, name)
	if err != nil {
		fmt.Println("updateNickname:", err)
		return utils.Err
	}
	return utils.Success
}

//update password
func Updatepwd(name string, pwd string) int {
	code := FindUser(name)
	if code == utils.Err {
		return utils.Err
	}
	if code == utils.ErrUserNotExit {
		return utils.ErrUserNotExit
	}
	p := utils.Md5(pwd)
	_, err := _updatePwd.Exec(p, name)
	if err != nil {
		return utils.Err
	}
	return utils.Success
}

//upload profile picture
func UploadPic(name string, picture string) int {
	code := FindUser(name)
	if code == utils.Err {
		return utils.Err
	}
	if code == utils.ErrUserNotExit {
		return utils.ErrUserNotExit
	}
	_, err := _uploadPic.Exec(picture, name)
	if err != nil {
		return utils.Err
	}
	return utils.Success
}

//get nickName and Profile picture
func GetInfo(name string) (nickname string, picture string, c int) {
	code := FindUser(name)
	if code == utils.Err {
		return "", "", utils.Err
	}
	if code == utils.ErrUserNotExit { //用户不存在
		return "", "", utils.ErrUserNotExit
	}
	rows, err := _getInfo.Query(name)
	if err != nil {
		fmt.Println("getInfo:", err)
		return nickname, picture, utils.Err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&nickname, &picture)
	}
	if err != nil {
		return "", "", utils.Err
	}
	return nickname, picture, utils.Success
}

//login
func Login(name string, password string) int {
	var pwd string
	code := FindUser(name)
	if code == utils.Err {
		return utils.Err
	}
	if code == utils.ErrUserNotExit { //用户不存在
		return utils.ErrUserNotExit
	}
	//fmt.Println(name, password)
	rows, err := _findUser.Query(name)
	if err != nil {
		fmt.Println("finduser:", err)
		return utils.Err
	}
	defer rows.Close()
	//get encrypted password
	for rows.Next() {
		err = rows.Scan(&pwd)
	}
	if err != nil {
		fmt.Println("login:", err)
		return utils.Err
	}
	//match
	p := utils.Md5(password)
	//p := password
	//not match
	if p != pwd {
		return utils.ErrPwdWrong
	}
	return utils.Success
}
