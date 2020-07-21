package schema

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

// UserSignByCredentialInput 账号密码登录
type UserSignByCredentialInput struct {
	Username string
	Password string
	Captcha  string
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
