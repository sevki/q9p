package main // import "sevki.org/q9p/q9pfs"

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
	"sevki.org/q9p/filesystem"
	"sevki.org/q9p/protocol"
)

const addr = "localhost:4242"

func main() {
	flag.Parse()

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}

	filesystemlistener, err := filesystem.Newfilesystem(func(l *protocol.Listener) error {
		l.Trace = nil // log.Printf
		return nil
	})

	if err := filesystemlistener.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}