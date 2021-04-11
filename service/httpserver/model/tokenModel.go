package model

import "time"

type VerifyCode struct {
	Email                string `json:"Email" bson:"email"`
	VerifyCode           string `json:"VerifyCode" bson:"verify_code"`
	VerifyCodeCreateTime string `json:"VerifyCodeCreateTime" bson:"verify_code_create_time"`
}

type AccessTokenCode struct {
	AccessToken           string    `json:"AccessToken" bson:"access_token"`
	UserId                string    `json:"UserId" bson:"user_id"`
	AccessTokenCreateTime time.Time `json:"AccessTokenCreateTime" bson:"access_token_create_date"`
	AccessTokenModifyTime time.Time `json:"AccessTokenModifyTime" bson:"access_token_modify_time"`
}
