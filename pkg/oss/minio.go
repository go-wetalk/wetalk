package oss

import (
	"appsrv/pkg/config"

	"github.com/minio/minio-go/v6"
)

var (
	Server = ""
)

var mc *minio.Client

// ProvideSingleton provides singleton instance of minio.Client.
func ProvideSingleton() *minio.Client {
	if mc == nil {
		var err error
		srvConf := config.ProvideSingleton()
		mc, err = minio.New(srvConf.Oss.Endpoint, srvConf.Oss.Access, srvConf.Oss.Secret, srvConf.Oss.Secure)
		if err != nil {
			panic(err)
		}

		exist, err := mc.BucketExists(srvConf.Oss.Bucket)
		if err != nil {
			panic(err)
		}

		if !exist {
			err = mc.MakeBucket(srvConf.Oss.Bucket, "")
			if err != nil {
				panic(err)
			}

			policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + srvConf.Oss.Bucket + `/*"],"Sid": ""}]}`
			err = mc.SetBucketPolicy(srvConf.Oss.Bucket, policy)
			if err != nil {
				panic(err)
			}
		}

		Server = mc.EndpointURL().String() + "/" + srvConf.Oss.Bucket
	}

	return mc
}
