FROM golang:1.16-alpine3.13 AS builder

RUN apk add --no-cache bash curl make wget
WORKDIR /go/src/asnlookup
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make && make install

FROM alpine:3.13
COPY --from=builder /usr/local/bin/asnlookup /usr/local/bin/asnlookup-utils /usr/local/bin/
USER nobody
ENV ASNLOOKUP_DB=/default.db
ENTRYPOINT ["/usr/local/bin/asnlookup"]
CMD []
