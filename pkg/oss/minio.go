package oss

import (
	"appsrv/pkg/config"

	"github.com/minio/minio-go/v6"
)

var (
	Bucket   = ""
	Endpoint = ""
	Server   = ""
)

var Min *minio.Client

// InitOss 获取对象存储操作实例
func InitOss(c config.OssConfig) (err error) {
	Min, err = minio.New(c.Endpoint, c.Access, c.Secret, c.Secure)
	if err != nil {
		return
	}

	exist, err := Min.BucketExists(c.Bucket)
	if err != nil {
		return
	}

	if !exist {
		err = Min.MakeBucket(c.Bucket, "")
		if err != nil {
			return
		}

		policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + c.Bucket + `/*"],"Sid": ""}]}`
		err = Min.SetBucketPolicy(c.Bucket, policy)
		if err != nil {
			return err
		}
	}

	Bucket = c.Bucket
	Endpoint = c.Endpoint
	Server = Min.EndpointURL().String() + "/" + Bucket

	return err
}
