package fileformat

import (
	"github.com/twinj/uuid"
	"path"
	"strings"
)

func UniqueFormat(fn string) string {
	// path.Ext() get the extension of the file
	fileName := strings.TrimSuffix(fn, path.Ext(fn))
	extension := path.Ext(fn)
	u := uuid.NewV4()
	newFileName := fileName + "-" + u.String() + extension

	return newFileName
}