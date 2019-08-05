package fda

import (
	"path/filepath"
	"strings"
)

// checking file extension.
type fileType struct {
	types      map[string]struct{}
	includeAll bool
	excludeDot bool
}

func newFileType(list []string) fileType {
	types := make(map[string]struct{})
	for _, s := range list {
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		types["."+s] = struct{}{}
	}

	includeAll := false
	if len(types) == 0 {
		includeAll = true
	}

	return fileType{
		types:      types,
		includeAll: includeAll,
		excludeDot: true,
	}
}

func (f *fileType) setIncludeAll(b bool) {
	f.includeAll = b
}

func (f *fileType) setExcludeDot(b bool) {
	f.excludeDot = b
}

func (f fileType) isTarget(path string) bool {
	if f.includeAll {
		if !f.excludeDot {
			return true
		}

		b := filepath.Base(path)
		if len(b) == 0 {
			return false
		}
		return string(b[0]) != "."
	}

	ext := strings.ToLower(filepath.Ext(path))
	_, ok := f.types[ext]
	return ok
}
