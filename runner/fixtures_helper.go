package runner

import (
	"path"
	"runtime"
)

func CurrentTestDir() string {
	_, filePath, _, _ := runtime.Caller(1)
	return path.Dir(filePath)
}
