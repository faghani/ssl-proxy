package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	m := autocert.Manager{
		Cache:  autocert.DirCache("./certs"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			return nil
		},
	}

	go func() {
		HTTPHandler := m.HTTPHandler(serveHTTP())
		http.ListenAndServe(":http", HTTPHandler)
	}()

	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   serveHTTP(),
	}

	log.Fatal(s.ListenAndServeTLS("", ""))
}

func serveHTTP() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		forwarder, _ := forward.New(forward.PassHostHeader(true))
		req.URL = testutils.ParseURI(os.Getenv("PROXY_TO"))
		res.Header().Set("Strict-Transport-Security", "max-age=86400; includeSubDomains")

		forwarder.ServeHTTP(res, req)
	})
}
