package logic

import (
	"github.com/swayrider/grpcclients/regionclient"
)

func newRegionClient(host string, port int) (*regionclient.Client, error) {
	clnt, err := regionclient.New(func() (string, int) { return host, port })
	if err != nil {
		return nil, err
	}
	return clnt.(*regionclient.Client), nil
}

func RegionSearchPoint(
	host string,
	port int,
	lat, lon float64,
	includeExtended bool,
) (regionclient.RegionList, error) {
	client, err := newRegionClient(host, port)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	defer client.Close()

	return client.SearchPoint(
		regionclient.Coordinate{Latitude: lat, Longitude: lon},
		includeExtended,
	)
}

func RegionSearchBox(
	host string,
	port int,
	minLat, minLon, maxLat, maxLon float64,
	includeExtended bool,
) (regionclient.RegionList, error) {
	client, err := newRegionClient(host, port)
	if err != nil {
		return regionclient.RegionList{}, err
	}
	defer client.Close()

	return client.SearchBox(
		regionclient.BoundingBox{
			BottomLeft: regionclient.Coordinate{Latitude: minLat, Longitude: minLon},
			TopRight:   regionclient.Coordinate{Latitude: maxLat, Longitude: maxLon},
		},
		includeExtended,
	)
}

func RegionFindPath(
	host string,
	port int,
	fromRegion, toRegion string,
) ([]string, error) {
	client, err := newRegionClient(host, port)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.FindRegionPath(fromRegion, toRegion)
}
