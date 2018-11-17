package examples

import (
	"errors"

	"github.com/mbecker/gpxs/geo"
)

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
	if previousPoint.Timestamp.Valid && p1.Timestamp.Valid {
		return p1.Timestamp.Time.Sub(*previousPoint.Timestamp.Time).Seconds(), nil
	}
	return 0, errors.New("Point or Previous Point does not have a timestamp")
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
	gpxPoint.Point.SetPointData(&previousGPXPoint.Point, algorithm)
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

// CheckActivityType returns the activity type (as a string number) based on my experience with strava, garmin, runkeeeper, ...
func (c *CustomAlgorithm) CheckActivityType(lowerCaseName string) (string, error) {

	var activityTpesName map[string]string
	activityTpesName = make(map[string]string)
	activityTpesName["running"] = "9"
	activityTpesName["lauf"] = "9"
	activityTpesName["cycling"] = "1"
	activityTpesName["rad"] = "1"
	activityTpesName["walking"] = "4"
	activityTpesName["hiking"] = "4"
	activityTpesName["spaziergang"] = "4"

	result := activityTpesName[lowerCaseName]
	if len(result) == 0 {
		return "", errors.New("Activity tpe not found")
	}
	return result, nil

}
