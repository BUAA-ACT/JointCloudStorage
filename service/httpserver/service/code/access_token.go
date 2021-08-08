package code

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	//"time"
	//"github.com/chilts/sid"
	//"github.com/kjk/betterguid"
	//"github.com/oklog/ulid"
	//"github.com/rs/xid"
	//"github.com/segmentio/ksuid"
	//"github.com/sony/sonyflake"
	"github.com/satori/go.uuid"
	//"github.com/dgrijalva/jwt-go"
)

//func genXid() {
//	id := xid.New()
//	fmt.Printf("github.com/rs/xid:           %s\n", id.String())
//}
//
//func genKsuid() {
//	id := ksuid.New()
//	fmt.Printf("github.com/segmentio/ksuid:  %s\n", id.String())
//}
//
//func genBetterGUID() {
//	id := betterguid.New()
//	fmt.Printf("github.com/kjk/betterguid:   %s\n", id)
//}
//
//func genUlid() {
//	t := time.Now().UTC()
//	entropy := rand.New(rand.NewSource(t.UnixNano()))
//	id := ulid.MustNew(ulid.Timestamp(t), entropy)
//	fmt.Printf("github.com/oklog/ulid:       %s\n", id.String())
//}
//
//func genSonyflake() {
//	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
//	id, err := flake.NextID()
//	if err != nil {
//		log.Fatalf("flake.NextID() failed with %s\n", err)
//	}
//	// Note: this is base16, could shorten by encoding as base62 string
//	fmt.Printf("github.com/sony/sonyflake:   %x\n", id)
//}
//
//func genSid() {
//	id := sid.Id()
//	fmt.Printf("github.com/chilts/sid:       %s\n", id)
//}

func genUUIDv4() uuid.UUID {
	id := uuid.NewV4()
	return id
}

func GenToken() uuid.UUID {
	return genUUIDv4()
}

//func createJWTToken(SecretKey []byte, issuer string, Uid uint, isAdmin bool) (tokenString string, err error) {
//	claims := jwt.CustomClaims{
//		jwt.StandardClaims{
//			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
//			Issuer: issuer,
//		},
//		Uid,
//		isAdmin,
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err = token.SignedString(SecretKey)
//	return
//}

// ------------------------------------------------------------

//type Claims struct {
//	UserId uint
//	UUID uuid.UUID
//	jwt.StandardClaims
//}
//
////登录以后签发jwt
//func tokenNext(c *gin.Context, user model.User) (string,error){
//	j := &middleware.JWT{
//		// 唯一签名
//		[]byte(global.GVA_CONFIG.JWT.SigningKey),
//	}
//	clams := request.CustomClaims{
//		UUID:        user.UUID,
//		ID:          user.UserId,
//		StandardClaims: jwt.StandardClaims{
//			//签名生效时间
//			NotBefore: int64(time.Now().Unix() - 1000),
//			//发放时间
//			IssuedAt:time.Now().Unix(),
//			//过期时间
//			ExpiresAt: int64(time.Now().Unix() + 60*60*24*7),
//			//签名的发行者
//			Issuer:    "ldy",
//		},
//	}
//	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
//	tokenString,err:=token.SignedString(jwtkey)
//	if err != nil {
//		return "",err
//	}
//	return tokenString,nil
//}
//
//====================================================================
//
////验证前端给的Token
//func AuthMiddleware()gin.HandlerFunc {
//
//	return func(context *gin.Context) {
//		//获取 header
//		tokenString := context.GetHeader("Authorization")
//
//		//验证格式
//		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
//			context.JSON(http.StatusUnauthorized, gin.H{
//				"code": 401,
//				"msg":  "权限不足",
//			})
//			context.Abort()
//			return
//		}
//
//		tokenString = tokenString[7:]
//		//验证tokenString
//		token, claims, err := common.ParseToken(tokenString)
//		if err != nil || !token.Valid {
//			context.JSON(http.StatusUnauthorized, gin.H{
//				"code": 401,
//				"msg":  "权限不足",
//			})
//			context.Abort()
//			return
//		}
//
//		//通过验证
//		userId := claims.UserId
//		DB := common.GetDB()
//		var user model.User
//		DB.First(&user, userId)
//
//		//用户不存在
//		if user.ID == 0 {
//			context.JSON(http.StatusUnauthorized, gin.H{
//				"code": 401,
//				"msg":  "权限不足",
//			})
//			context.Abort()
//			return
//		}
//		//将用户信息写入上写文
//		context.Set("user", user)
//
//		context.Next()
//	}
//}
//	//具体验证token方法
//func ParseToken(tokenStr string) (*jwt.Token,*Claims,error) {
//	claims:=&Claims{}
//	token,err:=jwt.ParseWithClaims(tokenStr,claims, func(token *jwt.Token) (i interface{}, err error) {
//		return jwt,nil
//	})
//	return token,claims,err
//}

//--------------------------------------

// Example (atypical) using the StandardClaims type by itself to parse a token.
// The StandardClaims type is designed to be embedded into your custom types
// to provide standard validation features.  You can use it alone, but there's
// no way to retrieve other fields after parsing.
// See the CustomClaimsType example for intended usage.
func ExampleNewWithClaims_standardClaims() {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.QsODzZu3lUZMVdhbO76u3Jv02iYCvEHcYVUI1kOWEU0 <nil>
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ExampleNewWithClaims_customClaimsType() {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c <nil>
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ExampleParseWithClaims_customClaimsType() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// sample token is expired.  override time so it parses as valid
	at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			fmt.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
		}
	})

	// Output: bar 15000
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	jwt.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwt.TimeFunc = time.Now
}

// An example of parsing the error types using bitfield checks
func ExampleParse_errorChecking() {
	// Token from another example.  This token is expired
	var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	// Output: Timing is everything
}
