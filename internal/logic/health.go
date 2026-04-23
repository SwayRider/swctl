package logic

import (
	"github.com/swayrider/grpcclients/healthclient"
)

func newHealthClient(host string, port int) (*healthclient.Client, error) {
	clnt, err := healthclient.New(func() (string, int) { return host, port })
	if err != nil {
		return nil, err
	}
	return clnt.(*healthclient.Client), nil
}

func Ping(host string, port int) error {
	client, err := newHealthClient(host, port)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Ping()
}

func HealthCheck(host string, port int, component string) (healthclient.ServiceStatus, error) {
	client, err := newHealthClient(host, port)
	if err != nil {
		return "", err
	}
	defer client.Close()
	return client.Check(component)
}
