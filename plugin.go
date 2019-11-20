package main

import (
	"io/ioutil"
	pathutil "path"
	"strings"
	"time"

	"github.com/yingce/drone-oss-cache/cachekey"

	log "github.com/sirupsen/logrus"
	"github.com/yingce/drone-oss-cache/lib/cache/archive/util"
	"github.com/yingce/drone-oss-cache/lib/cache/cache"
	"github.com/yingce/drone-oss-cache/lib/cache/storage"
)

// Plugin structure
type Plugin struct {
	Filename     string
	Path         string
	FallbackPath string
	FlushPath    string
	Mode         string
	FlushAge     int
	Mount        []string
	Cacert       string
	CacertPath   string

	Storage storage.Storage
}

const (
	// RestoreMode for resotre mode string
	RestoreMode = "restore"
	// RebuildMode for rebuild mode string
	RebuildMode = "rebuild"
	// FlushMode for flush mode string
	FlushMode = "flush"
)

// Exec runs the plugin
func (p *Plugin) Exec() error {
	var err error

	var useCheckSum bool

	if strings.Contains(p.Path, "checksum") || strings.Contains(p.Filename, "checksum") {
		useCheckSum = true
	}

	p.Path, err = cachekey.CacheKey(p.Path, cachekey.MetaData{})
	if err != nil {
		log.Fatal(err)
	}
	p.Filename, err = cachekey.CacheKey(p.Filename, cachekey.MetaData{})
	if err != nil {
		log.Fatal(err)
	}
	p.FallbackPath, err = cachekey.CacheKey(p.FallbackPath, cachekey.MetaData{})
	if err != nil {
		log.Fatal(err)
	}

	at, err := util.FromFilename(p.Filename)

	if err != nil {
		return err
	}

	c := cache.New(p.Storage, at)

	path := pathutil.Join(p.Path, p.Filename)
	fallbackPath := pathutil.Join(p.FallbackPath, p.Filename)

	if p.Cacert != "" {
		certPath := "/etc/ssl/certs/ca-certificates.crt"
		log.Infof("Installing new ca certificate at %s", certPath)
		err := installCaCert(certPath, p.Cacert)

		if err == nil {
			log.Info("Successfully installed new certificate")
		}
	}

	if p.CacertPath != "" {
		certPath := "/etc/ssl/certs/ca-certificates.crt"
		log.Infof("Installing new ca certificate at %s", certPath)
		err := installCaCertFromPath(certPath, p.CacertPath)

		if err == nil {
			log.Info("Successfully installed new certificate")
		}
	}

	if p.Mode == RebuildMode {
		log.Infof("Rebuilding cache at %s", path)
		var exists bool
		if useCheckSum {
			exists, _ = p.Storage.Exists(path)
		}
		if !exists {
			err = c.Rebuild(p.Mount, path)
			if err == nil {
				log.Infof("Cache rebuilt")
			}
		} else {
			log.Infof("Cache skip, object exists[using checksum func]")
		}
	}

	if p.Mode == RestoreMode {
		log.Infof("Restoring cache at %s", path)
		err = c.Restore(path, fallbackPath)

		if err == nil {
			log.Info("Cache restored")
		}
	}

	if p.Mode == FlushMode {
		log.Infof("Flushing cache items older than %d days at %s", p.FlushAge, path)
		f := cache.NewFlusher(p.Storage, genIsExpired(p.FlushAge))
		err = f.Flush(p.FlushPath)

		if err == nil {
			log.Info("Cache flushed")
		}
	}

	return err
}

func genIsExpired(age int) cache.DirtyFunc {
	return func(file storage.FileEntry) bool {
		// Check if older than "age" days
		return file.LastModified.Before(time.Now().AddDate(0, 0, age*-1))
	}
}

func installCaCert(path, cacert string) error {
	err := ioutil.WriteFile(path, []byte(cacert), 0644)
	return err
}

func installCaCertFromPath(path, cacertPath string) error {
	cacert, err := ioutil.ReadFile(cacertPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(cacert), 0644)
	return err
}
