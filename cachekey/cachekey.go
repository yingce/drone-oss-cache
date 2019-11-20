package cachekey

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"text/template"
	"time"
)

type MetaData map[string]interface{}

func CacheKey(path string, data MetaData) (string, error) {
	tmpl, err := template.New("cachePath").Funcs(funcMap).Parse(path)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

var funcMap = template.FuncMap{
	"checksum": func(path string) string {
		absPath, err := filepath.Abs(filepath.Clean(path))
		if err != nil {
			log.Println("cache key template/checksum could not find file")
			return ""
		}

		f, err := os.Open(absPath)
		if err != nil {
			log.Println("cache key template/checksum could not open file")
			return ""
		}
		defer f.Close()

		str, err := readerHasher(f)
		if err != nil {
			log.Println("cache key template/checksum could not generate hash")
			return ""
		}
		return str
	},
	"epoch": func() string { return strconv.FormatInt(time.Now().Unix(), 10) },
	"arch":  func() string { return runtime.GOARCH },
	"os":    func() string { return runtime.GOOS },
}

func readerHasher(readers ...io.Reader) (string, error) {
	h := md5.New() // #nosec

	for _, r := range readers {
		if _, err := io.Copy(h, r); err != nil {
			return "", fmt.Errorf("write reader as hash %w", err)
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
