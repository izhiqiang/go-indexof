package logger

import (
	"fmt"
	"os"
)

func Logger(a ...any) {
	fmt.Fprint(os.Stdout, a...)
}
