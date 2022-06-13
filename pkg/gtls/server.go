package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type Server struct {
	CaFile   string
	CertFile string
	KeyFile  string
}

func (s *Server) GetCredentialsByCA() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(s.CaFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("certPool.AppendCertsFromPEM err")
	}
	credential := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return credential, err
}

func (s *Server) GetTLSCredentials() (credentials.TransportCredentials, error) {
	c, err := credentials.NewServerTLSFromFile(s.CertFile, s.KeyFile)
	if err != nil {
		return nil, err
	}

	return c, err
}
