package assets

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"simple-gateway/config"
)

func NewAssetsFileServer(assets config.Assets) (string, http.Handler, error) {
	source, err := url.Parse(assets.Source)
	if err != nil {
		return "", nil, err
	}
	stat, err := os.Stat(assets.Dir)
	if err != nil {
		return "", nil, err
	}

	if !stat.IsDir() {
		return "", nil, fmt.Errorf("需要一个文件目录，而不是文件")
	}

	return source.Host, http.StripPrefix(assets.StripPrefix, http.FileServer(http.Dir(assets.Dir))), nil
}
