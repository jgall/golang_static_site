FROM golang as goBuilder
WORKDIR /go/src
COPY main.go .

RUN set -x 
RUN go get -d -v .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o serve 

# Runtime Stage
FROM alpine
WORKDIR /app
COPY --from=goBuilder /go/src/serve .
COPY --from=goBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./serve"]