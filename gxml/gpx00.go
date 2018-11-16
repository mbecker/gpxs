package gxml

import (
	"encoding/xml"

	"github.com/mbecker/gpxs/generic"
)

// Defines all base xml elements used for v1.0 and v1.1

//GPX00GpxLink struct fields for all links
type GPX00GpxLink struct {
	Href string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

//GPX00GpxPoint struct fields for all points
type GPX00GpxPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	// Position info
	Ele         generic.NullableFloat64 `xml:"ele,omitempty"`
	Timestamp   string                  `xml:"time,omitempty"`
	MagVar      string                  `xml:"magvar,omitempty"`
	GeoIdHeight string                  `xml:"geoidheight,omitempty"`
	// Description info
	Name  string         `xml:"name,omitempty"`
	Cmt   string         `xml:"cmt,omitempty"`
	Desc  string         `xml:"desc,omitempty"`
	Src   string         `xml:"src,omitempty"`
	Links []GPX00GpxLink `xml:"link"`
	Sym   string         `xml:"sym,omitempty"`
	Type  string         `xml:"type,omitempty"`
	// Accuracy info
	Fix           string   `xml:"fix,omitempty"`
	Sat           *int     `xml:"sat,omitempty"`
	Hdop          *float64 `xml:"hdop,omitempty"`
	Vdop          *float64 `xml:"vdop,omitempty"`
	Pdop          *float64 `xml:"pdop,omitempty"`
	AgeOfDGpsData *float64 `xml:"ageofdgpsdata,omitempty"`
	DGpsID        *int     `xml:"dgpsid,omitempty"`
}

//GPX00GpxRte struct fields for a route
type GPX00GpxRte struct {
	XMLName xml.Name `xml:"rte"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO:
	//Links       []Link   `xml:"link"`
	Number generic.NullableInt `xml:"number,omitempty"`
	Type   string              `xml:"type,omitempty"`
	Points []*GPX00GpxPoint    `xml:"rtept"`
}

//GPX10GpxTrk struct fiels for a track
type GPX00GpxTrk struct {
	XMLName xml.Name `xml:"trk"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO:
	//Links    []Link   `xml:"link"`
	Number   generic.NullableInt `xml:"number,omitempty"`
	Type     string              `xml:"type,omitempty"`
	Segments []*GPX00GpxTrkSeg   `xml:"trkseg,omitempty"`
}

//GPX00GpxTrkSeg strcut fields for all track segements
type GPX00GpxTrkSeg struct {
	XMLName xml.Name         `xml:"trkseg"`
	Points  []*GPX00GpxPoint `xml:"trkpt"`
}
