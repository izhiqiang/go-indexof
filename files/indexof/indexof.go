package indexof

import (
	"fmt"
	"indexof/config"
	"net/http"
	"os"
	"time"
)

var (
	Global IIndexOf

	IndexMap = map[string]IIndexOf{
		"local": &Local{},
	}
)

func LoadIndexOf() error {
	name := config.Global.IndexOf.Name
	iof, ok := IndexMap[name]
	if !ok {
		return fmt.Errorf("load config not find name: %s", name)
	}
	Global = iof
	Global.SetRoot(config.Global.IndexOf.Root)
	return nil
}

type IIndexOf interface {
	SetRoot(root string)
	ReadURLPath(indexOf string) ([]PathInfo, error)
	FindURLPath(URLPath string) (os.FileInfo, error)
	URLPathDownload(URLPath string, w http.ResponseWriter, r *http.Request)
}

type PathInfo struct {
	IsDir   bool      `json:"is_dir"`
	IsImage bool      `json:"is_image"`
	Name    string    `json:"name"`
	Size    string    `json:"size"`
	Data    time.Time `json:"data"`
}
