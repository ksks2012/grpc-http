# Reference

https://github.com/eddycjy
# intruduction

An example for grpc to provide http interface

# install protobuf compiler


## download & install

```
wget https://github.com/google/protobuf/releases/download/v3.11.2/protobuf-all-3.11.2.zip
unzip protobuf-all-3.11.2.zip && cd protobuf-3.11.2/
./configure
make
make install
```
## check version

```
ldconfig
protoc --version
```

## go grpc install

```
go get -u google.golang.org/grpc
```

## Protoc Plugin

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

# compile proto

## basic
```
protoc --go_out=plugins=grpc:. ./proto/*.proto
```

## with google.api.http
```
protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-scosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. ./proto/*.proto
```

## gen code with buf

### reference

https://docs.buf.build/tour/introduction

### gen buf.yaml

```bash
buf config init
```

### Update mode

```bash
buf mod update
```

### gen go code

```bash
buf generate
```


# build project
```
go build ./cmd/grpc-http
```

# execute
```
./grpc-http
```

# grpcurl

## install
```
go get github.com/fullstorydev/grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

## test
```
grpcurl -plaintext localhost:8001 proto.TagService.GetTagList
```


# cmux

## install
```
go get -u github.com/soheilhy/cmux@v0.1.4
```

# grpc-gateway

## install
```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

## move binary file
```
mv $GOPATH/bin/protoc-gen-grpc-gateway /usr/local/go/bin/
```
