package examples

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
)

func example1() {
	// 1.) Use a built-in geo.Algorithm
	vincenty := geo.Vincenty{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               1.644854, // ~95%
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:                    1 / 298.257223563,
		SemiMinorAxisB:                6356752.314245,
		Epsilon:                       1e-12,
		MaxIterations:                 200,
	}

	// 2.) Parse a gpx file with the geo.Algorithm
	gpxDoc, err := gpxs.ParseFile(filepath.Join("fileDirectory", "test.gpx"), &vincenty)
	if err != nil {
		panic(err)
	}

	var (
		distance        float64
		duration        float64
		movingDistance  float64
		movingTime      float64
		stoppedDistance float64
		stoppedTime     float64

		trackDistance        float64
		trackDuration        float64
		trackMovingDistance  float64
		trackMovingTime      float64
		trackStoppedDistance float64
		trackStoppedTime     float64

		segmentDistance        float64
		segmentDuration        float64
		segmentMovingDistance  float64
		segmentMovingTime      float64
		segmentStoppedDistance float64
		segmentStoppedTime     float64
	)

	// GPX
	md := gpxDoc.MovingData
	distance += md.Distance
	duration += md.Duration
	movingDistance += md.MovingDistance
	movingTime += md.MovingTime
	stoppedDistance += md.StoppedDistance
	stoppedTime += md.StoppedTime

	// Tracks
	for _, track := range gpxDoc.Tracks {
		trackMd := track.MovingData
		trackDistance += trackMd.Distance
		trackDuration += trackMd.Duration
		trackMovingDistance += trackMd.MovingDistance
		trackMovingTime += trackMd.MovingTime
		trackStoppedDistance += trackMd.StoppedDistance
		trackStoppedTime += trackMd.StoppedTime

		// Segments
		for _, segment := range track.Segments {
			segmentMd := segment.MovingData
			segmentDistance += segmentMd.Distance
			segmentDuration += segmentMd.Duration
			segmentMovingDistance += segmentMd.MovingDistance
			segmentMovingTime += segmentMd.MovingTime
			segmentStoppedDistance += segmentMd.StoppedDistance
			segmentStoppedTime += segmentMd.StoppedTime
		}
	}

	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(duration)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(movingTime)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(stoppedTime)))

	t04, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackDuration)))
	t05, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackMovingTime)))
	t06, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackStoppedTime)))

	t07, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentDuration)))
	t08, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentMovingTime)))
	t09, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentStoppedTime)))

	fmt.Printf("------")
	fmt.Printf("%s", t01)                    // GPX Duration
	fmt.Printf("%f", distance/1000.0)        // GPX Distance
	fmt.Printf("%s", t02)                    // GPX Moving Time)
	fmt.Printf("%s", t03)                    // GPX Stopped Time
	fmt.Printf("%f", movingDistance/1000.0)  // GPX Moving Distance
	fmt.Printf("%f", stoppedDistance/1000.0) // GPX Stopped Distance
	fmt.Printf("------")
	fmt.Printf("%s", t04)                         // Track Duration
	fmt.Printf("%f", trackDistance/1000.0)        // Track Moving Distance
	fmt.Printf("%s", t05)                         // Track Moving Time
	fmt.Printf("%s", t06)                         // Track Stopped Time
	fmt.Printf("%f", trackMovingDistance/1000.0)  // Track Moving Distance
	fmt.Printf("%f", trackStoppedDistance/1000.0) // Track Stopped Distance
	fmt.Printf("------")
	fmt.Printf("%s", t07)                           // Segment Duration
	fmt.Printf("%f", segmentDistance/1000.0)        // Segment Distance
	fmt.Printf("%s", t08)                           // Segment Moving Time
	fmt.Printf("%s", t09)                           // Segment  Stopped Time
	fmt.Printf("%f", segmentMovingDistance/1000.0)  // Segment Moving Distance
	fmt.Printf("%f", segmentStoppedDistance/1000.0) // Segment  Stopped Distance
}
