package indexof

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Local struct {
	root string
}

func (l *Local) SetRoot(root string) {
	l.root = root
}

func (l *Local) FindURLPath(URLPath string) (os.FileInfo, error) {
	URLPath = filepath.Join(l.root, URLPath)
	fileInfo, err := os.Stat(URLPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist", URLPath)
	}
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func (l *Local) ReadURLPath(URLPath string) ([]PathInfo, error) {
	readDir, err := os.ReadDir(filepath.Join(l.root, URLPath))
	if err != nil {
		return nil, err
	}
	var pis []PathInfo
	for _, entry := range readDir {
		info, _ := entry.Info()
		pis = append(pis, PathInfo{
			IsDir: entry.IsDir(),
			Name:  entry.Name(),
			Size:  autoConvertSize(info.Size()),
			Data:  info.ModTime(),
		})
	}
	return pis, err
}

// URLPathDownload 文件下载
// curl -I http://127.0.0.1:8080/filename
// curl -r 0-60 http://127.0.0.1:8080/filename -o part1
// curl -r 60-81 http://127.0.0.1:8080/filename -o part1
// cat part1 part2 > filename
func (l *Local) URLPathDownload(URLPath string, w http.ResponseWriter, r *http.Request) {
	URLPath = filepath.Join(l.root, URLPath)
	file, err := os.Open(URLPath)
	info, err2 := file.Stat()
	if err != nil || err2 != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	fileName := info.Name()
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to stat files", http.StatusInternalServerError)
		return
	}
	fileSize := stat.Size()
	rangeHeader := r.Header.Get("Range")
	w.Header().Set("Accept-Ranges", "bytes")
	if rangeHeader != "" {
		rangeParts := strings.Split(rangeHeader, "=")
		if len(rangeParts) == 2 {
			rangeValue := rangeParts[1]
			rangeParts = strings.Split(rangeValue, "-")
			if len(rangeParts) == 2 {
				start, _ := strconv.ParseInt(rangeParts[0], 10, 64)
				end := fileSize - 1
				if rangeParts[1] != "" {
					end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
				}
				if start > end || start >= fileSize {
					http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
					return
				}
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
				w.Header().Set("Accept-Ranges", "bytes")
				w.WriteHeader(http.StatusPartialContent)
				_, err = file.Seek(start, 0)
				if err != nil {
					http.Error(w, "Unable to seek files", http.StatusInternalServerError)
					return
				}
				_, _ = io.CopyN(w, file, end-start+1)
				return
			}
		}
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	http.ServeFile(w, r, URLPath)
}

// AutoConvertSize 自动转换文件大小
func autoConvertSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	if size < KB {
		return fmt.Sprintf("%d B", size)
	} else if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	} else if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
}
