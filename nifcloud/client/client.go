package client

import (
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devopsrunner"
	"github.com/nifcloud/nifcloud-sdk-go/service/dns"
	"github.com/nifcloud/nifcloud-sdk-go/service/ess"
	"github.com/nifcloud/nifcloud-sdk-go/service/nas"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
	"github.com/nifcloud/nifcloud-sdk-go/service/storage"
)

// Client nifcloud client
type Client struct {
	Computing    *computing.Client
	RDB          *rdb.Client
	NAS          *nas.Client
	DNS          *dns.Client
	ESS          *ess.Client
	Storage      *storage.Client
	DevOps       *devops.Client
	DevOpsRunner *devopsrunner.Client
}

// New return Client
func New(cfg nifcloud.Config, storageCfg nifcloud.Config) *Client {
	return &Client{
		Computing:    computing.NewFromConfig(cfg),
		RDB:          rdb.NewFromConfig(cfg),
		NAS:          nas.NewFromConfig(cfg),
		DNS:          dns.NewFromConfig(cfg),
		ESS:          ess.NewFromConfig(cfg),
		Storage:      storage.NewFromConfig(storageCfg),
		DevOps:       devops.NewFromConfig(cfg),
		DevOpsRunner: devopsrunner.NewFromConfig(cfg),
	}
}
