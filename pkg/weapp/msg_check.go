// Package weapp 提供微信服务器端接口调用整合
package weapp

import (
	"errors"

	"github.com/medivhzhan/weapp/v2"
)

var ErrMsgSec = errors.New("您的输入包含违规内容哟")

// MsgCheck 检查文本是否有违规内容
func MsgCheck(s string) error {
	t, err := GetAccessToken()
	if err != nil {
		return err
	}

	res, err := weapp.MSGSecCheck(t, s)
	if err != nil {
		return err
	}

	if err = res.GetResponseError(); err != nil {
		return ErrMsgSec
	}

	return nil
}
