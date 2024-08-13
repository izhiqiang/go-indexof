package handle

import (
	"fmt"
	"indexof/files/indexof"
	"indexof/view"
	"net/http"
	"runtime"
)

type IndexOfResponse struct {
	IndexOf   string
	GoVersion string
	PathInfos []indexof.PathInfo
}

func IndexOfHandler(w http.ResponseWriter, r *http.Request) {
	resp := IndexOfResponse{IndexOf: r.URL.Path, GoVersion: runtime.Version(), PathInfos: []indexof.PathInfo{}}
	fileInfo, err := indexof.Global.FindURLPath(resp.IndexOf)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", resp.IndexOf), http.StatusNotFound)
		return
	}
	if !fileInfo.IsDir() {
		indexof.Global.URLPathDownload(resp.IndexOf, w, r)
		return
	}
	resp.PathInfos, _ = indexof.Global.ReadURLPath(resp.IndexOf)
	_ = view.FetchIndexOf(w, resp)
}
