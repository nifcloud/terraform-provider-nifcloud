package client

import (
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/rdb"
)

// Client nifcloud client
type Client struct {
	Computing *computing.Client
	RDB       *rdb.Client
}

// New return Client
func New(cfg nifcloud.Config) *Client {
	return &Client{
		Computing: computing.New(cfg),
		RDB:       rdb.New(cfg),
	}
}
