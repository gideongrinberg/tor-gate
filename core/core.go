/*
	This package contains the core logic for the proxy.
*/

package core

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/gideongrinberg/tor-gate/assets"
	"golang.org/x/net/proxy"
)

var config *Config
var logger *log.Logger

func Logger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return logger
}

func createTorClient() *http.Client {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 dialer: %v", err)
	}

	dialCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext:       dialCtx,
			DisableKeepAlives: true,
		},

		Timeout: 30 * time.Second, // TODO: good value?
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if config.ShowDisclaimer {
		cookie, err := r.Cookie("disclaimer_acknowledged")
		if err != nil || cookie.Value != "true" {
			w.Header().Add("Content-Type", "text/html")
			w.Write(assets.Disclaimer)
			return
		}
	}

	host := r.Host
	if colon := strings.Index(host, ":"); colon != -1 {
		host = host[:colon]
	}

	// TODO: handle invalid subdomains
	parts := strings.Split(host, ".")
	subdomain := strings.Join(parts[:len(parts)-2], ".")
	if config.EnableTranslations {
		onion, err := config.Translations[subdomain]
		if !err {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Could not find %s.onion", subdomain)))
		}

		subdomain = onion
	}

	targetOnion := "http://" + subdomain + ".onion" + r.RequestURI

	if slices.Contains(config.Blacklist, subdomain) {
		w.WriteHeader(403)
		w.Write(assets.Blacklist)
		return
	}

	if config.WhitelistOnly && !slices.Contains(config.Whitelist, subdomain) {
		w.WriteHeader(403)
		w.Write(assets.Whitelist)
		return
	}

	client := createTorClient()
	body, _ := io.ReadAll(r.Body) // need error handling
	req, _ := http.NewRequest("GET", targetOnion, bytes.NewReader(body))
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to proxy request through Tor."))
		log.Println(err)
		return
	}

	defer res.Body.Close()

	for name, values := range res.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(res.StatusCode)

	if strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		html := RewriteLinks(buf.Bytes(), config.Domain)
		w.Write([]byte(html))
	} else {
		io.Copy(w, res.Body)
	}

}

func StartServer() {
	config = LoadConfig()
	s := &http.Server{
		Addr:           config.Port,
		Handler:        http.HandlerFunc(handleRequest),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		Logger().Printf("Starting server on %s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger().Fatalf("Could not listen on %s: %v\n", s.Addr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // wait for sigint
	Logger().Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		Logger().Printf("Forcibly shutting down due to error: %v", err)
	}

	Logger().Println("Successfully terminated server")
}
