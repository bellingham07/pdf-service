package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQL struct {
		DSN             string
		LogMode         bool
		MaxOpenCons     int
		MaxIdleCons     int
		CreateBatchSize int
	}
}
