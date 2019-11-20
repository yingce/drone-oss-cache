package util

import (
	"fmt"
	"strings"

	"github.com/yingce/drone-oss-cache/lib/cache/archive"
	"github.com/yingce/drone-oss-cache/lib/cache/archive/tar"
	"github.com/yingce/drone-oss-cache/lib/cache/archive/tgz"
)

// FromFilename determines the archive format to use based on the name.
func FromFilename(name string) (archive.Archive, error) {
	if strings.HasSuffix(name, ".tar") {
		return tar.New(), nil
	}

	if strings.HasSuffix(name, ".tgz") || strings.HasSuffix(name, ".tar.gz") {
		return tgz.New(), nil
	}

	return nil, fmt.Errorf("Unknown file format for archive %s", name)
}
