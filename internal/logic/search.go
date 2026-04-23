package logic

import (
	"fmt"

	"github.com/swayrider/grpcclients/searchclient"
)

type SearchResult struct {
	label       string
	locality    string
	region      string
	country     string
	confidence  float64
	layer       string
	lat         float64
	lon         float64
	street      string
	houseNumber string
	id          string
	localAdmin  string
	countryCode string
	name        string
}

func (r *SearchResult) Label() string       { return r.label }
func (r *SearchResult) Locality() string    { return r.locality }
func (r *SearchResult) Region() string      { return r.region }
func (r *SearchResult) Country() string     { return r.country }
func (r *SearchResult) Confidence() float64 { return r.confidence }
func (r *SearchResult) Layer() string       { return r.layer }
func (r *SearchResult) Lat() float64        { return r.lat }
func (r *SearchResult) Lon() float64        { return r.lon }
func (r *SearchResult) Street() string      { return r.street }
func (r *SearchResult) HouseNumber() string { return r.houseNumber }
func (r *SearchResult) Id() string          { return r.id }
func (r *SearchResult) LocalAdmin() string  { return r.localAdmin }
func (r *SearchResult) CountryCode() string { return r.countryCode }
func (r *SearchResult) Name() string        { return r.name }

func (r *SearchResult) Display() {
	fmt.Printf("\t%s (%s, %s) [%s] lat=%.6f lon=%.6f confidence=%.2f\n",
		r.label, r.locality, r.country, r.layer, r.lat, r.lon, r.confidence)
}

func newSearchResult(
	label, locality, region, country string,
	confidence float64,
	layer string,
	lat, lon float64,
	street, housenumber, id, localadmin, countryCode, name string,
) searchclient.SearchResult {
	return &SearchResult{
		label: label, locality: locality, region: region, country: country,
		confidence: confidence, layer: layer, lat: lat, lon: lon,
		street: street, houseNumber: housenumber, id: id, localAdmin: localadmin,
		countryCode: countryCode, name: name,
	}
}

func newSearchClient(host string, port int) (*searchclient.Client, error) {
	clnt, err := searchclient.New(func() (string, int) { return host, port })
	if err != nil {
		return nil, err
	}
	return clnt.(*searchclient.Client), nil
}

func Geocode(
	authHost string,
	authPort int,
	searchHost string,
	searchPort int,
	user string,
	password string,
	query string,
	size int32,
	language string,
	focusLat float64,
	focusLon float64,
) (results []*SearchResult, err error) {
	authClient, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer authClient.Close()

	accessToken, _, err := authClient.Login(user, password, false)
	if err != nil {
		return
	}

	client, err := newSearchClient(searchHost, searchPort)
	if err != nil {
		return
	}
	defer client.Close()

	q := searchclient.SearchQuery{
		Text: query,
		Viewport: searchclient.BoundingBox{
			BottomLeft: searchclient.Coordinate{Latitude: -90, Longitude: -180},
			TopRight:   searchclient.Coordinate{Latitude: 90, Longitude: 180},
		},
		Size:     size,
		Language: language,
	}
	if focusLat != 0 || focusLon != 0 {
		q.FocusPoint = &searchclient.Coordinate{Latitude: focusLat, Longitude: focusLon}
	}

	l, err := client.Search(accessToken, q, newSearchResult)
	if err != nil {
		return
	}

	results = make([]*SearchResult, 0, len(l))
	for _, r := range l {
		results = append(results, r.(*SearchResult))
	}
	return
}

func ReverseGeocode(
	authHost string,
	authPort int,
	searchHost string,
	searchPort int,
	user string,
	password string,
	lat float64,
	lon float64,
	size int32,
	language string,
) (results []*SearchResult, err error) {
	authClient, err := newAuthClient(authHost, authPort)
	if err != nil {
		return
	}
	defer authClient.Close()

	accessToken, _, err := authClient.Login(user, password, false)
	if err != nil {
		return
	}

	client, err := newSearchClient(searchHost, searchPort)
	if err != nil {
		return
	}
	defer client.Close()

	q := searchclient.ReverseGeocodeQuery{
		Point:    searchclient.Coordinate{Latitude: lat, Longitude: lon},
		Size:     size,
		Language: language,
	}

	l, err := client.ReverseGeocode(accessToken, q, newSearchResult)
	if err != nil {
		return
	}

	results = make([]*SearchResult, 0, len(l))
	for _, r := range l {
		results = append(results, r.(*SearchResult))
	}
	return
}

var _ searchclient.SearchResult = (*SearchResult)(nil)
