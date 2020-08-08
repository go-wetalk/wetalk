package schema

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID   uint
	Name string
	Logo string
}

// UserLoginByWeappInput 微信小程序登录参数
type UserLoginByWeappInput struct {
	Code          string `validate:"required" json:"code"`
	EncryptedData string `validate:"required" json:"encrypted_data"`
	IV            string `validate:"required" json:"iv"`
	RawData       string `validate:"required" json:"raw_data"`
	Signature     string `validate:"required" json:"signature"`
	UserInfo      struct {
		AvatarURL string `json:"avatarUrl"`
		NickName  string `json:"nickName"`
		Gender    int    `json:"gender"`
	} `validate:"required" json:"user_info"`
}

// UserLoginByQappInput QQ小程序登录
type UserLoginByQappInput struct {
	Code          string `validate:"required" json:"code"`
	EncryptedData string `validate:"required" json:"encrypted_data"`
	IV            string `validate:"required" json:"iv"`
	RawData       string `validate:"required" json:"raw_data"`
	Signature     string `validate:"required" json:"signature"`
	UserInfo      struct {
		AvatarURL string `json:"avatarUrl"`
		NickName  string `json:"nickName"`
		Gender    int    `json:"gender"`
	} `validate:"required" json:"user_info"`
}

// UserSignUpInput 账号密码登录
type UserSignUpInput struct {
	Username string
	Password string
	Captcha  string
}

func (v UserSignUpInput) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Username, validation.Required.Error("请输入用户名"), validation.RuneLength(1, 12).Error("用户名限制12个字符")),
		validation.Field(&v.Password, validation.Required.Error("请设置密码"), validation.RuneLength(6, 32).Error("密码长度限制6到32个字符")),
	)
}

type UserSignOutput struct {
	Token string
	User  struct {
		ID        uint
		Name      string
		AvatarURL string
		Gender    int
		Created   string
		Coin      int
	}
}

type UserStatus struct {
	ID           uint
	Name         string
	Logo         string
	Gender       int
	Coin         int
	Created      string
	RoleList     []string
	UnreadNotify int
}

type UserDetail struct {
	User

	Created      string
	TopicCount   int
	Topics       []TopicListItem
	CommentCount int
	Comments     []CommentBadge
}
