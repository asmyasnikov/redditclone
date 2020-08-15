package redditclone

import (
	"fmt"
	"github.com/asmyasnikov/redditclone/api"
)

func Run(cfg *api.Configuration) error {
	fmt.Println(cfg.Server.HTTPListen)
	return nil
}
