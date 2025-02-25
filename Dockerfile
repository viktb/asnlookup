FROM golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache bash curl make wget
WORKDIR /go/src/asnlookup
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make && make install

FROM alpine:3.21
COPY --from=builder /usr/local/bin/asnlookup /usr/local/bin/asnlookup-utils /usr/local/bin/
USER nobody
ENV ASNLOOKUP_DB=/default.db
ENTRYPOINT ["/usr/local/bin/asnlookup"]
CMD []
