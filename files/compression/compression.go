package compression

import "io"

type Compression interface {
	SetIoWriter(w io.Writer)
	Do(path string) error
}
