package qapp

// 检测地址
const (
	apiIMGSecCheck = "/api/json/security/ImgSecCheck"
	apiMSGSecCheck = "/api/json/security/MsgSecCheck"
)

// IMGSecCheck 本地图片检测
// 官方文档: https://developers.weixin.qq.com/miniprogram/dev/api/imgSecCheck.html
//
// filename 要检测的图片本地路径
// token 接口调用凭证(access_token)
func IMGSecCheck(token, filename string) (*CommonError, error) {
	api := baseURL + apiIMGSecCheck
	return imgSecCheck(api, token, filename)
}

func imgSecCheck(api, token, filename string) (*CommonError, error) {

	url, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	res := new(CommonError)
	if err := postFormByFile(url, "media", filename, res); err != nil {
		return nil, err
	}

	return res, nil
}

func IMGSecCheck2(token, url string) (*CommonError, error) {
	api := baseURL + apiIMGSecCheck
	return imgSecCheck2(api, token, url)
}

func imgSecCheck2(api, token, fileURL string) (*CommonError, error) {
	url, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	res := new(CommonError)
	if err := postFormByFileURL(url, "media", fileURL, res); err != nil {
		return nil, err
	}

	return res, nil
}

// MSGSecCheck 文本检测
// 官方文档: https://developers.weixin.qq.com/miniprogram/dev/api/msgSecCheck.html
//
// content 要检测的文本内容，长度不超过 500KB，编码格式为utf-8
// token 接口调用凭证(access_token)
func MSGSecCheck(token, content string) (*CommonError, error) {
	api := baseURL + apiMSGSecCheck
	return msgSecCheck(api, token, content)
}

func msgSecCheck(api, token, content string) (*CommonError, error) {
	url, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	params := requestParams{
		"content": content,
	}

	res := new(CommonError)
	if err = postJSON(url, params, res); err != nil {
		return nil, err
	}

	return res, nil
}
