package aliyun_oss

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/yingce/drone-oss-cache/lib/cache/storage"
)

// Options contains configuration for the S3 connection.
type Options struct {
	Endpoint string
	Key      string
	Secret   string
}

type ossStorage struct {
	client *oss.Client
	opts   *Options
}

// New method creates an implementation of Storage with S3 as the backend.
func New(opts *Options) (storage.Storage, error) {
	client, err := oss.New(opts.Endpoint, opts.Key, opts.Secret)
	if err != nil {
		return nil, err
	}
	return &ossStorage{
		client: client,
		opts:   opts,
	}, nil
}

func (s *ossStorage) Get(p string, dst io.Writer) error {
	bucket, key := splitBucket(p)

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("Invalid path %s", p)
	}

	log.Infof("Retrieving file in %s at %s", bucket, key)

	exists, err := s.client.IsBucketExist(bucket)

	if !exists {
		log.Infof("Bucket %s already exists", bucket)
		return err
	}

	bkt, err := s.client.Bucket(bucket)

	if err != nil {
		return err
	}

	object, err := bkt.GetObject(key)

	if err != nil {
		return err
	}

	log.Infof("Copying object from the server")

	numBytes, err := io.Copy(dst, object)

	if err != nil {
		return err
	}

	log.Infof("Downloaded %s from server", humanize.Bytes(uint64(numBytes)))

	return nil
}

func (s *ossStorage) findOrCreateBucket(bucketName string) (bucket *oss.Bucket, err error) {
	exists, err := s.client.IsBucketExist(bucketName)
	if err != nil {
		return
	}

	if !exists {
		if err = s.client.CreateBucket(bucketName); err != nil {
			return
		}
		log.Infof("Bucket %s created", bucketName)
	} else {
		log.Infof("Bucket %s already exists", bucketName)
	}
	bucket, err = s.client.Bucket(bucketName)
	if err != nil {
		return
	}
	return
}

func (s *ossStorage) Put(p string, src io.Reader) error {
	bucket, key := splitBucket(p)

	log.Infof("Uploading to bucket %s at %s", bucket, key)

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("Invalid path %s", p)
	}

	bkt, err := s.findOrCreateBucket(bucket)

	if err != nil {
		return err
	}

	log.Infof("Putting file in %s at %s", bucket, key)

	options := []oss.Option{oss.ContentType("application/tar")}

	err = bkt.PutObject(key, src, options...)

	if err != nil {
		return err
	}

	h, err := bkt.GetObjectMeta(key)

	if err != nil {
		return err
	}

	length := h.Get("Content-Length")
	lengthInt, _ := strconv.Atoi(length)

	log.Infof("Uploaded %s to server", humanize.Bytes(uint64(lengthInt)))

	return nil
}

func (s *ossStorage) List(p string) ([]storage.FileEntry, error) {
	bucket, key := splitBucket(p)

	log.Infof("Retrieving objects in bucket %s at %s", bucket, key)

	if len(bucket) == 0 || len(key) == 0 {
		return nil, fmt.Errorf("Invalid path %s", p)
	}

	exists, err := s.client.IsBucketExist(bucket)

	if err != nil {
		return nil, fmt.Errorf("%s does not exist: %s", p, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s does not exist", p)
	}

	bkt, err := s.client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	var objects []storage.FileEntry
	hasPrefix := strings.HasSuffix(key, "/")
	if !hasPrefix {
		key += "/"
	}
	res, err := bkt.ListObjects(oss.Prefix(key))
	if err != nil {
		return nil, err
	}
	for _, object := range res.Objects {
		path := bucket + "/" + object.Key
		objects = append(objects, storage.FileEntry{
			Path:         path,
			Size:         object.Size,
			LastModified: object.LastModified,
		})
		log.Debugf("Found object %s: Path=%s Size=%d LastModified=%s", object.Key, path, object.Size, object.LastModified)
	}

	log.Infof("Found %d objects in bucket %s at %s", len(objects), bucket, key)

	return objects, nil
}

func (s *ossStorage) Exists(p string) (bool, error) {
	bucket, key := splitBucket(p)
	bkt, err := s.findOrCreateBucket(bucket)
	if err != nil {
		return false, nil
	}
	return bkt.IsObjectExist(key)
}

func (s *ossStorage) Delete(p string) error {
	bucket, key := splitBucket(p)

	log.Infof("Deleting object in bucket %s at %s", bucket, key)

	if len(bucket) == 0 || len(key) == 0 {
		return fmt.Errorf("Invalid path %s", p)
	}

	exists, err := s.client.IsBucketExist(bucket)

	if err != nil {
		return fmt.Errorf("%s does not exist: %s", p, err)
	}
	if !exists {
		return fmt.Errorf("%s does not exist", p)
	}

	bkt, err := s.client.Bucket(bucket)
	if err != nil {
		return err
	}

	err = bkt.DeleteObject(key)
	return err
}

func splitBucket(p string) (string, string) {
	// Remove initial forward slash
	full := strings.TrimPrefix(p, "/")

	// Get first index
	i := strings.Index(full, "/")

	if i != -1 && len(full) != i+1 {
		// Bucket names need to be all lower case for the key it doesnt matter
		return strings.ToLower(full[0:i]), full[i+1:]
	}

	return "", ""
}
