package dao

import (
	"fmt"
	"time"

	"GoUserManaSys/utils"

	"github.com/go-redis/redis"
)

var _c *redis.Client

func init() {
	_c = redis.NewClient(&redis.Options{
		Addr:     utils.RedisAddress,
		Password: "",
		DB:       0,
		PoolSize: utils.RedisPoolSize,
	})
	_, err := _c.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("init redis success...")
}

//Close
func CloseRedis() {
	if err := _c.Close(); err != nil {
		fmt.Println("redis closed failed")
	}
	fmt.Println("redis closed...")

}

//redis get the user information
func RdsGetInfo(name string) (nickName string, pic string, hasData bool, errCode int) {
	v, err := _c.HGetAll(name).Result()
	if err != nil {
		return "", "", false, utils.ErrRedisGet
	}
	//valid初始为空，如果不为空，则证明有效
	if v["valid"] != "" {
		hasData = true
	} else {
		hasData = false
	}
	return v["nick_name"], v["profile_picture"], hasData, utils.Success
}

//redis set user information
func RdsSetInfo(name string, nickname string, pic string, expTime int64) int {
	v := map[string]interface{}{
		"valid":           "1",
		"nick_name":       nickname,
		"profile_picture": pic,
	}
	if err := _c.HMSet(name, v).Err(); err != nil {
		return utils.ErrRedisSet
	}
	_c.Expire(name, time.Duration(expTime*1e9))
	return utils.Success
}

//delete user information
func RdsDelInfo(name string) {
	_c.Del(name)
}

//set user invalid
func Invalid(name string) int {
	if err := _c.HSet(name, "valid", "").Err(); err != nil {
		return utils.ErrRedisSet
	}
	return utils.Success
}

//add token
func SetToken(name string, token string, expTime int64) int {
	err := _c.Set("auth_"+name, token, time.Duration(expTime*1e9)).Err()
	if err != nil {
		return utils.ErrRedisSet
	}
	return utils.Success
}

//check whether the tokens matched
func CheckToken(name string, token string) int {
	//stress test token whether match
	if token == utils.TestToken {
		return utils.Success
	}
	//release token whether match
	v, err := _c.Get("auth_" + name).Result()
	if err != nil {
		//获取token出错
		return utils.ErrRedisGet
	}
	if token == v {
		//token匹配
		return utils.Success
	}
	//token don't match
	return utils.ErrTokenWrong
}

//set token -> username
func SetTokenName(token string, name string, expTime int64) int {
	err := _c.Set(token, name, time.Duration(expTime*1e9)).Err()
	if err != nil {
		return utils.ErrRedisSet
	}
	return utils.Success
}

//get token -> username
func GetTokenName(token string) (name string, err error) {
	name, err = _c.Get(token).Result()
	if err != nil {
		return "", err
	}
	return name, nil
}
