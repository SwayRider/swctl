package logic

import (
	"context"

	"github.com/swayrider/grpcclients/regionclient"
)

func newRegionClient(host string, port int) (*regionclient.Client, error) {
	clnt, err := regionclient.New(func() (string, int) { return host, port })
	if err != nil {
		return nil, err
	}
	return clnt.(*regionclient.Client), nil
}

func regionToken(authHost string, authPort int, user, password string) (string, error) {
	authClient, err := newAuthClient(authHost, authPort)
	if err != nil {
		return "", err
	}
	defer authClient.Close()
	token, _, err := authClient.Login(user, password, false)
	return token, err
}

func RegionSearchPoint(
	authHost string,
	authPort int,
	user, password string,
	host string,
	port int,
	lat, lon float64,
	includeExtended bool,
) (regionclient.RegionList, error) {
	token, err := regionToken(authHost, authPort, user, password)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	client, err := newRegionClient(host, port)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	defer client.Close()

	return client.SearchPoint(
		context.Background(),
		token,
		regionclient.Coordinate{Latitude: lat, Longitude: lon},
		includeExtended,
	)
}

func RegionSearchBox(
	authHost string,
	authPort int,
	user, password string,
	host string,
	port int,
	minLat, minLon, maxLat, maxLon float64,
	includeExtended bool,
) (regionclient.RegionList, error) {
	token, err := regionToken(authHost, authPort, user, password)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	client, err := newRegionClient(host, port)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	defer client.Close()

	return client.SearchBox(
		context.Background(),
		token,
		regionclient.BoundingBox{
			BottomLeft: regionclient.Coordinate{Latitude: minLat, Longitude: minLon},
			TopRight:   regionclient.Coordinate{Latitude: maxLat, Longitude: maxLon},
		},
		includeExtended,
	)
}

func RegionFindPath(
	authHost string,
	authPort int,
	user, password string,
	host string,
	port int,
	fromRegion, toRegion string,
) ([]string, error) {
	token, err := regionToken(authHost, authPort, user, password)
	if err != nil {
		return nil, err
	}
	client, err := newRegionClient(host, port)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.FindRegionPath(context.Background(), token, fromRegion, toRegion)
}
