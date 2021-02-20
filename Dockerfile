FROM  golang  AS builder
ADD . /go/src/chainmeta-exporter/
RUN cd /go/src/chainmeta-exporter && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o chain-meta-exporter

From alpine
COPY --from=builder /go/src/chainmeta-exporter/chain-meta-exporter /usr/bin/chain-meta-exporter
CMD ["/usr/bin/chain-meta-exporter"]   

