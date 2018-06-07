# easyhttps
A golang package for easy integration of https (using TLS) into existing http servers

### Usage:

    func main() {
      srv := &http.Server{
        Handler: http.FileServer(http.Dir("static")),
        Addr:    ":80",
      }

      httpsSrv := easyhttps.WrapHTTPS(srv, ":443", "tlsCertCache", "example.com")
      log.Fatal(httpsSrv.ListenAndServe())
    }
   