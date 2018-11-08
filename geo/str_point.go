package geo

import (
	"time"

	"github.com/mbecker/gpxs/generic"
)

// Point represents generic point data and implements the Location interface
type Point struct {
	Latitude  float64
	Longitude float64
	Elevation generic.NullableFloat64
	Timestamp time.Time

	Distance float64 // The distance (m) from the previos to this point
	Duration float64 // The duration (sec) from the previos to this point
	Speed    float64 // The speed (m/s) from the previous point to this point
	Pace     float64 // The pace (m/s) from the previous point to this point
}

//GetLatitude returns the latitude
func (pt *Point) GetLatitude() float64 {
	return pt.Latitude
}

//GetLongitude returns the longititude
func (pt *Point) GetLongitude() float64 {
	return pt.Longitude
}

//GetElevation returns the elevation
func (pt *Point) GetElevation() generic.NullableFloat64 {
	return pt.Elevation
}

// GetDistance returns the distance of two GpxWpts with the given algorithm.Distance
// func (pt *Point) GetDistance(loc2 Location, algorithm Algorithm) (float64, error) {
// 	return algorithm.Distance(pt, &loc2)
// }

// SetPointData sets the the point data for duration, distance, speed, pace
func (pt *Point) SetPointData(prevPoint Point, algorithm Algorithm) {
	// Duration (sec)
	duration, errDuration := algorithm.Duration(pt, &prevPoint)
	if errDuration != nil {
		duration = 0
	}
	pt.Duration = duration

	// Distance (m)
	distance, errDistance := algorithm.Distance(pt, &prevPoint)
	if errDistance != nil {
		distance = 0
	}
	pt.Distance = distance

	// Speed (m/s)
	speed, errSpeed := algorithm.Speed(distance, duration)
	if errSpeed != nil {
		speed = 0
	}
	pt.Speed = speed

	// Pace
	pace, errPace := algorithm.Pace(distance, duration)
	if errPace != nil {
		pace = 0
	}
	pt.Pace = pace
}
