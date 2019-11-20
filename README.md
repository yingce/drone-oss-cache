# drone-s3-cache

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-s3-cache/status.svg)](http://cloud.drone.io/drone-plugins/drone-s3-cache)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/yingce/drone-oss-cache.svg)](https://microbadger.com/images/yingce/drone-oss-cache "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-s3-cache?status.svg)](http://godoc.org/github.com/drone-plugins/drone-s3-cache)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-s3-cache)](https://goreportcard.com/report/github.com/drone-plugins/drone-s3-cache)

Drone plugin that allows you to cache directories within the build workspace, this plugin is backed by S3 compatible storages. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-s3-cache/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-oss-cache
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag yingce/drone-oss-cache .
```

## Usage

Support providers: S3[default], OSS  

```console
docker run --rm \
  -e PLUGIN_FLUSH=true \
  -e PLUGIN_ENDPOINT="http://minio.company.com" \
  -e PLUGIN_ACCESS_KEY="myaccesskey" \
  -e PLUGIN_SECRET_KEY="mysecretKey" \
  -e PLUGIN_FILENAME="{{checksum \"go.mod\"}}" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  yingce/drone-oss-cache

Aliyun OSS provider
docker run --rm \
  -e PLUGIN_PROVIDER=OSS \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT="https://oss-cn-beijing.aliyuncs.com" \
  -e PLUGIN_ACCESS_KEY="myaccesskey" \
  -e PLUGIN_SECRET_KEY="mysecretKey" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  yingce/drone-oss-cache

docker run --rm \
  -e PLUGIN_RESTORE=true \
  -e PLUGIN_ENDPOINT="http://minio.company.com" \
  -e PLUGIN_ACCESS_KEY="myaccesskey" \
  -e PLUGIN_SECRET_KEY="mysecretKey" \
  -e DRONE_REPO_OWNER="foo" \
  -e DRONE_REPO_NAME="bar" \
  -e DRONE_COMMIT_BRANCH="test" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  yingce/drone-oss-cache

docker run -it --rm \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  alpine:latest sh -c "mkdir -p cache && echo 'testing cache' >> cache/test && cat cache/test"

docker run --rm \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_MOUNT=".bundler" \
  -e PLUGIN_ENDPOINT="http://minio.company.com" \
  -e PLUGIN_ACCESS_KEY="myaccesskey" \
  -e PLUGIN_SECRET_KEY="mysecretKey" \
  -e DRONE_REPO_OWNER="foo" \
  -e DRONE_REPO_NAME="bar" \
  -e DRONE_COMMIT_BRANCH="test" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  yingce/drone-oss-cache
```
