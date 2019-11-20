package tgz

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"compress/gzip"
	"io"

	"github.com/yingce/drone-oss-cache/lib/cache/archive"
	"github.com/yingce/drone-oss-cache/lib/cache/archive/tar"
)

type tgzArchive struct{}

// New creates an archive that uses the .tar.gz file format.
func New() archive.Archive {
	return &tgzArchive{}
}

func (a *tgzArchive) Pack(srcs []string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	taP := tar.New()

	err := taP.Pack(srcs, gw)

	return err
}

func (a *tgzArchive) Unpack(dst string, r io.Reader) error {
	gr, err := gzip.NewReader(r)

	if err != nil {
		return err
	}

	taU := tar.New()

	fwErr := taU.Unpack(dst, gr)

	return fwErr
}
