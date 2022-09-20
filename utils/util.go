package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/google/uuid"
)

const (
	TestToken = "testToken"
)

//use md5 algorithm to encrypt password
func Md5(pwd string) string {
	p := md5.New()
	p.Write([]byte(pwd))
	return hex.EncodeToString(p.Sum(nil))
}

//use username to get token
//func GetToken(name string) (token string) {
//	s := md5.New()
//	s.Write([]byte(name + strconv.FormatInt(time.Now().Unix(), 10)))
//	token = hex.EncodeToString(s.Sum(nil))
//	return token
//}

//use uuid as token
func GetToken() string {
	uuid := uuid.New()
	key := uuid.String()
	return key
}

//test token
func GetTokenTest() string {
	return TestToken
}

//check image type valid?
func CheckImage(name string) bool {
	n := path.Ext(name)
	if n == ".png" || n == ".jpg" || n == ".gif" || n == ".jpeg" {
		return true
	}
	return false
}

//generate new name of image to avoid repeat
func NewImgName(name string) string {
	n := path.Ext(name)
	//生成uuid
	uuid := uuid.New()
	key := uuid.String()
	return key + n
}

//copy image
func SetImg(file multipart.File, head *multipart.FileHeader) error {
	newName := NewImgName(head.Filename)
	filepath := StaticFilePath + newName
	showFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	defer showFile.Close()
	_, err = io.Copy(showFile, file)
	if err != nil {
		return err
	}
	return nil
}
