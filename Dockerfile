FROM golang:1.10 as builder
WORKDIR /go/src/github.com/tsg/gotpl
RUN go get -d -v gopkg.in/yaml.v2
COPY tpl.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix gotpl -o gotpl .

FROM alpine
COPY --from=builder /go/src/github.com/tsg/gotpl/gotpl /usr/local/bin/gotpl
ENTRYPOINT ["/usr/local/bin/gotpl"]
