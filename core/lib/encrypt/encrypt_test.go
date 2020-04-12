package encrypt

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	md5, err := GetMD5([]byte("hahahahah"))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(HexEncodeToString(md5))
	}
}
func TestRandom(t *testing.T) {
	str, er := GetRandomSalt([]byte(""), 10)
	t.Log(HexEncodeToString(str), er)
}

func TestJwt(t *testing.T) {
	//	写入 JWT
	//	10分钟后过期
	timeDuration := 2 * time.Second

	//	要写入 JWT 的加密参数
	mapClaims := jwt.MapClaims{
		"UserId":         3322,
		"UserName":       "hahahah",
		"ExpirationDate": time.Now().Add(timeDuration).Format("2006-01-02 15:04:05"),
	}
	//	如果加密失败，则返回错误
	if tokenString, err := JWTEncrypt(mapClaims); err != nil {
		t.Error(err)
	} else {
		t.Log(tokenString)
		time.Sleep(5 * time.Second)
		if st, e := JWTDecrypt(tokenString); e != nil {
			t.Error(err)
		} else {
			t.Log(st)
		}
	}
}
