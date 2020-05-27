module appsrv

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-pg/pg/v9 v9.1.6
	github.com/go-pg/urlstruct v0.4.0 // indirect
	github.com/go-redis/cache/v7 v7.0.2
	github.com/go-redis/redis/v7 v7.2.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/jinzhu/now v1.1.1
	github.com/kataras/muxie v1.0.9
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/laeo/qapp v0.0.0-20200502160927-5198b1938446
	github.com/medivhzhan/weapp/v2 v2.1.1
	github.com/minio/minio-go/v6 v6.0.55
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/encoding v0.1.12 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cast v1.3.1
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.11
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	golang.org/x/image v0.0.0-20200430140353-33d19683fad8
	golang.org/x/net v0.0.0-20200519113804-d87ec0cfa476 // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/ini.v1 v1.56.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/medivhzhan/weapp/v2 => github.com/laeo/weapp/v2 v2.0.3-0.20200412002818-e35fbd5456cf
