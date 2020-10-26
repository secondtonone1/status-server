//go module proxy 配置  https://goproxy.cn
module status-server

go 1.13

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.8

replace golang.org/x/text => github.com/golang/text v0.3.3

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20200619210111-0f592d2728bb

replace golang.org/x/net => github.com/golang/net v0.0.0-20200602114024-627f9648deb9

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200604202706-70a84ac30bf9

replace golang.org/x/exp => github.com/golang/exp v0.0.0-20200513190911-00229845015e

replace cloud.google.com/go => github.com/googleapis/google-cloud-go v0.58.0

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20200620020550-bd6e04640131

replace google.golang.org/api => github.com/googleapis/google-api-go-client v0.28.0

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

replace go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/google/btree v1.0.0 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/micro/go-micro/v2 v2.9.0
	github.com/micro/go-plugins/registry/zookeeper/v2 v2.8.0
	github.com/pkg/errors v0.9.1
	github.com/tebeka/strftime v0.1.4 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	google.golang.org/grpc v1.29.1
)
