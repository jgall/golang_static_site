FROM golang as goBuilder
WORKDIR /go/src/github.com/jgall/golang_static_site
COPY . .

RUN set -x 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /serve 

# Runtime Stage
FROM alpine
WORKDIR /app
COPY --from=goBuilder /serve .
COPY --from=goBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["./serve"]