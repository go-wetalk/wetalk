package weapp

import (
	"appsrv/pkg/oss"
	"bytes"
	"io/ioutil"
	"time"

	"github.com/medivhzhan/weapp/v2"
	"github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
)

// GenerateWxaCode 生成微信二维码
func GenerateWxaCode(path string) (string, error) {
	token, err := GetAccessToken()
	if err != nil {
		return "", err
	}

	h := weapp.QRCode{Path: path}
	resp, res, err := h.Get(token)
	if err != nil {
		return "", err
	}

	if err = res.GetResponseError(); err != nil {
		return "", err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	objectName := time.Now().Format("20060102") + "/" + uuid.NewV4().String() + ".jpg"

	_, err = oss.Min.PutObject(
		oss.Bucket,
		objectName,
		bytes.NewReader(content),
		int64(len(content)),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		})
	if err != nil {
		return "", err
	}

	return objectName, nil
}
