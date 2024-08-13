package compression

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Zip struct {
	w io.Writer
}

func (z *Zip) SetIoWriter(w io.Writer) {
	z.w = w
}

func (z *Zip) Do(srcDir string) error {
	zipWriter := zip.NewWriter(z.w)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)
	err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过软连接
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		// 创建 ZIP 文件条目
		relPath, _ := filepath.Rel(srcDir, path)
		if info.IsDir() {
			relPath += "/"
		}
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// 复制文件内容到 ZIP 条目
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}
