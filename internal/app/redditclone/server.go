package redditclone

import (
	"fmt"
	"github.com/asmyasnikov/redditclone/api"
)

// Run run server with configuration cfg
func Run(cfg *api.Configuration) error {
	fmt.Println(cfg.Server.HTTPListen)
	return nil
}
