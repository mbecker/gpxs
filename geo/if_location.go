package geo

import "github.com/mbecker/gpxs/generic"

//Location implements an interface for all kinds of lat/long/elevation information
type Location interface {
	GetLatitude() float64
	GetLongitude() float64
	GetElevation() generic.NullableFloat64
}
