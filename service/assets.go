package service

import (
	"fmt"
	"github.com/qmstar0/nightsky-gateway/config"
	"net/http"
	"os"
)

func NewAssetsFileServer(assets *config.Assets) (http.Handler, error) {
	stat, err := os.Stat(assets.Dir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("需要一个文件目录，而不是文件")
	}

	return http.StripPrefix(assets.StripPrefix, http.FileServer(http.Dir(assets.Dir))), nil
}
