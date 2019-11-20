FROM 1.13 as builder

# Add Maintainer Info
LABEL maintainer="yingce@live.com"

ENV GO111MODULE on
# Set the GOPROXY environment variable
ENV GOPROXY https://goproxy.io

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN BUILD_OPTS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o drone-oss-cache

######## Start a new stage from scratch or alpine #######
FROM 3.10

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/drone-oss-cache .

# Command to run the executable
CMD ["/root/drone-oss-cache"]
