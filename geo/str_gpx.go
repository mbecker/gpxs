package geo

import (
	"fmt"
	"time"

	"github.com/mbecker/gpxs/generic"
)

//GPX implements one or multiple GPS tracks that can be written to and parsed
//from a gpx file
type GPX struct {
	XMLNs        string
	XmlNsXsi     string
	XmlSchemaLoc string

	Version          string
	Creator          string
	Name             string
	Description      string
	AuthorName       string
	AuthorEmail      string
	AuthorLink       string
	AuthorLinkText   string
	AuthorLinkType   string
	Copyright        string
	CopyrightYear    string
	CopyrightLicense string
	Link             string
	LinkText         string
	LinkType         string
	Timestamp        *time.Time
	Keywords         string

	MovingData MovingData

	// TODO:
	//Extensions []byte
	Routes    []GPXRoute
	Tracks    []GPXTrack
	Waypoints []GPXPoint
}

func (gpx *GPX) String() string {
	var result string
	result = fmt.Sprintf("--- GPX ---\n")
	result += fmt.Sprintf("Version: %s\n", gpx.Version)
	result += fmt.Sprintf("Name: %s\n", gpx.Name)
	result += fmt.Sprintf("Creator: %s\n", gpx.Creator)
	result += fmt.Sprintf("Time: %s\n", gpx.Timestamp)
	result += gpx.MovingData.String("--->")
	result += fmt.Sprintf("------\n")
	return result
}

//GPXRoute implements a gpx route
type GPXRoute struct {
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO:
	//Links       []Link
	Number int // generic.NullableInt
	Type   string
	// TODO:
	Points []GPXPoint
}

//GPXTrack implements a gpx track
type GPXTrack struct {
	Name        string
	Comment     string
	Description string
	Source      string
	Timestamp   *time.Time
	/**
	 * TODO:
	 * - [x] Should Links be included in tracks due to gpx specification?
	 */
	//Links    []Link
	Number     int //generic.NullableInt
	Type       string
	Segments   []GPXTrackSegment
	MovingData MovingData
}

func (track *GPXTrack) String() string {
	var result string
	result = fmt.Sprintf("--- Track ---\n")
	result += fmt.Sprintf("Name: %s\n", track.Name)
	result += fmt.Sprintf("Number: %v\n", track.Number)
	result += fmt.Sprintf("Time: %s\n", track.Timestamp)
	result += fmt.Sprintf("Segments #: %d\n", len(track.Segments))
	result += track.MovingData.String("::::")
	result += fmt.Sprintf("------\n")
	return result
}

//GPXTrackSegment represents a segment of a track
type GPXTrackSegment struct {
	Points     []GPXPoint
	MovingData MovingData
	Duration   float64
	Distance   float64
	Timestamp  *time.Time
	// TODO: extensions
}

func (seg *GPXTrackSegment) String() string {
	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(seg.Duration)))

	var result string
	result = fmt.Sprintf("--- GPXTrackSegment ---\n")
	result += fmt.Sprintf("Time: %s\n", seg.Timestamp)
	result += fmt.Sprintf("Duration Time: %s\n", t01)
	result += fmt.Sprintf("Distance: %f km\n", seg.Distance/1000.0)
	result += fmt.Sprintf("------\n")
	return result
}

//GPXPoint represents a point of the gpx file
type GPXPoint struct {
	Point
	// TODO:

	// TODO: Type
	MagneticVariation string
	// TODO: Type
	GeoidHeight string
	// Description info
	Name        string
	Comment     string
	Description string
	Source      string
	// TODO:
	// Links       []GpxLink
	Symbol string
	Type   string
	// Accuracy info
	TypeOfGpsFix       string
	Satellites         generic.NullableInt
	HorizontalDilution generic.NullableFloat64
	VerticalDilution   generic.NullableFloat64
	PositionalDilution generic.NullableFloat64
	AgeOfDGpsData      generic.NullableFloat64
	DGpsId             generic.NullableInt
}

//GpxBounds contains min/max latitude and longitude
type GpxBounds struct {
	MinLatitude  float64
	MaxLatitude  float64
	MinLongitude float64
	MaxLongitude float64
}
