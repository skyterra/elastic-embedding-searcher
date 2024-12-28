module github.com/skyterra/elastic-embedding-searcher

go 1.22.3

toolchain go1.22.10

require (
	github.com/bytedance/mockey v1.2.13
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/elastic/go-elasticsearch/v7 v7.17.10
	github.com/google/uuid v1.6.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/skyterra/clog v1.0.3
	github.com/spf13/cobra v1.8.1
	golang.org/x/net v0.32.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.0
)

require (
	github.com/gopherjs/gopherjs v1.12.80 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/arch v0.11.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)

replace (
	github.com/skyterra/clog v1.0.3 => ../clog
)