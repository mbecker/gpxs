package examples

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
)

// 1.) Define customer geo.algorithm; must fullfill interface geo.Algorithm
type CustomAlgorithm struct {
	CustomParameter float64
}

// ShouldStandardDeviation (CustomAlgorithm) returns if the standard deviation should be used or not
func (c *CustomAlgorithm) ShouldStandardDeviation() bool {
	return false
}

// Sigma (CustomAlgorithm) returns the sigma for the standard deviation; not used because of 'ShouldStandardDeviation' returns 'false'
func (c *CustomAlgorithm) Sigma() float64 {
	return 0
}

// Duration (CustomAlgorithm) returns the time.Duration from point p1 to previousPoint in sec
func (c *CustomAlgorithm) Duration(p1 *geo.Point, previousPoint *geo.Point) (float64, error) {
	return p1.Timestamp.Sub(previousPoint.Timestamp).Seconds(), nil
}

// CustomMovingPoints (CustomAlgorithm) defines which points should be used for "Moving"Time/Distance and if the it's set the new gpxPoint.Point Data
func (c *CustomAlgorithm) CustomMovingPoints(gpxPoint *geo.GPXPoint, previousGPXPoint *geo.GPXPoint, algorithm geo.Algorithm) error {

	/* 	Define which points should be used; if a point should be used for calculation then set it's new values like Duration, Distance, Speed, etc.
	Here we use the set the new value for the points which used for "Moving"Time/Distanc
	*/

	// speed < 100 m/s
	if gpxPoint.Speed < 100.0 {
		return errors.New("Point Speed below threshold")
	}
	gpxPoint.Point.SetPointData(previousGPXPoint.Point, algorithm)
	return nil
}

// Distance (CustomAlgorithm) returns just 100 as an example
func (c *CustomAlgorithm) Distance(p1 *geo.Point, previousPoint *geo.Point) (float64, error) {
	return 101, nil
}

// Speed (CustomAlgorithm) returns the speed in m/s
func (c *CustomAlgorithm) Speed(distance float64, duration float64) (float64, error) {
	return 101.9, nil
}

// Pace (CustomAlgorithm) returns the pace in s/m
func (c *CustomAlgorithm) Pace(distance float64, duration float64) (float64, error) {
	return 20.9, nil
}

func example2() {

	customAlgorithm := CustomAlgorithm{
		CustomParameter: 100.9,
	}

	// 2.) Parse a gpx file with the custom geo.Algorithm
	gpxDoc, err := gpxs.ParseFile(filepath.Join("fileDirectory", "test.gpx"), &customAlgorithm)
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
