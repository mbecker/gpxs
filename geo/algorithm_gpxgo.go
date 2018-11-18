package geo

import (
	"errors"
	"math"
	"strings"
)

//AlgorithmGpxgo defines the basic calculation of the distance (2D/3D) and the Haversine formula
type AlgorithmGpxgo struct {
	Name                          string
	ShouldStandardDeviationBeUsed bool // Should the standard deviation be used to determine which points are used for calculation
	SigmaMultiplier               float64
	ShouldHaversine               bool // Should the formula of Haversine be used to calculate the distance between two points
	OneDegree                     float64
	EarthRadius                   float64
	Should3D                      bool // Should the distance be calculated with the elevation different of each point
}

// String returns the name of the algorithm
func (alg *AlgorithmGpxgo) String() string {
	return alg.Name
}

// ShouldStandardDeviation (AlgorithmGpxgo) returns if the standard deviation should be used or not
func (alg *AlgorithmGpxgo) ShouldStandardDeviation() bool {
	return alg.ShouldStandardDeviationBeUsed
}

// Sigma (AlgorithmGpxgo) returns the sigma for the standard deviation
func (alg *AlgorithmGpxgo) Sigma() float64 {
	return alg.SigmaMultiplier
}

// Duration (AlgorithmGpxgo) returns the time.Duration from point p1 to previousPoint in sec
func (alg *AlgorithmGpxgo) Duration(p1 *Point, previousPoint *Point) (float64, error) {
	if previousPoint.Timestamp.Valid && p1.Timestamp.Valid {
		return p1.Timestamp.Time.Sub(*previousPoint.Timestamp.Time).Seconds(), nil
	}
	return 0, errors.New("Point or Previous Point does not have a timestamp")
}

// CustomMovingPoints (AlgorithmGpxgo) defines which points should be used for "Moving"Time/Distance and if the it's set the new gpxPoint.Point Data
func (alg *AlgorithmGpxgo) CustomMovingPoints(gpxPoint *GPXPoint, previousGPXPoint *GPXPoint, algorithm Algorithm) error {

	/* 	Define which points should be used; if a point should be used for calculation then set it's new values like Duration, Distance, Speed, etc.
	Here we use the set the new value for the points which used for "Moving"Time/Distanc
	*/

	// speed < 1 m/s
	if gpxPoint.Speed < 1.0 {
		return errors.New("Point Speed below threshold")
	}
	gpxPoint.Point.SetPointData(&previousGPXPoint.Point, algorithm)
	return nil
}

// Distance (AlgorithmGpxgo) returns either 2d or 3d distance or the length by the formula Haversine
func (alg *AlgorithmGpxgo) Distance(p1 *Point, previousPoint *Point) (float64, error) {

	absLat := math.Abs(p1.Latitude - previousPoint.Latitude)
	absLon := math.Abs(p1.Longitude - previousPoint.Longitude)
	if alg.ShouldHaversine || absLat > 0.2 || absLon > 0.2 {
		return HaversineDistance(p1.Latitude, p1.Longitude, previousPoint.Latitude, previousPoint.Longitude, alg.EarthRadius), nil
	}

	coef := math.Cos(ToRad(p1.Latitude))
	x := p1.Latitude - previousPoint.Latitude
	y := (p1.Longitude - previousPoint.Longitude) * coef

	distance2d := math.Sqrt(x*x+y*y) * alg.OneDegree

	if !alg.Should3D || p1.Elevation == previousPoint.Elevation {
		return distance2d, nil
	}

	eleDiff := 0.0
	if p1.Elevation.NotNull() && previousPoint.Elevation.NotNull() {
		eleDiff = p1.Elevation.Value() - previousPoint.Elevation.Value()
	}

	return math.Sqrt(math.Pow(distance2d, 2) + math.Pow(eleDiff, 2)), nil
}

// Speed (Vincenty) returns the speed in m/s
func (alg *AlgorithmGpxgo) Speed(distance float64, duration float64) (float64, error) {
	if duration == 0 {
		return 0, errors.New("Duration is zero")
	}
	speed := distance / duration
	if math.IsInf(speed, 1) {
		return 0, errors.New("Duration is +Inf")
	}
	if math.IsNaN(speed) {
		return 0, errors.New("Duration IsNaN")
	}

	return distance / duration, nil
}

// Pace (Vincenty) returns the pace in s/m
func (alg *AlgorithmGpxgo) Pace(distance float64, duration float64) (float64, error) {
	if math.IsInf(distance, 1) || math.IsInf(distance, -1) || math.IsNaN(distance) || math.IsInf(duration, 1) || math.IsNaN(duration) {
		return 0, errors.New("Distance is +INf or NaN")
	}
	if distance == 0 {
		return 0, errors.New("Distance is zero")
	}

	pace := duration / distance
	if math.IsInf(pace, 1) {
		return 0, errors.New("Duration is +Inf")
	}
	if math.IsNaN(pace) {
		return 0, errors.New("Duration IsNaN")
	}

	return pace, nil
}

/* Standard Algorithm internal methods */

//ToRad converts to radial coordinates
func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

// HaversineDistance returns the haversine distance between two points.
//
// Implemented from http://www.movable-type.co.uk/scripts/latlong.html
func HaversineDistance(lat1, lon1, lat2, lon2 float64, earthRadius float64) float64 {
	dLat := ToRad(lat1 - lat2)
	dLon := ToRad(lon1 - lon2)
	thisLat1 := ToRad(lat1)
	thisLat2 := ToRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c

	return d
}

// CheckActivityType returns the activity type (as a string number) based on my experience with strava, garmin, runkeeeper, ...
func (alg *AlgorithmGpxgo) CheckActivityType(lowerCaseName string) (string, error) {

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
	if len(result) > 0 {
		return result, nil
	}

	for key, value := range activityTpesName {
		if strings.Contains(lowerCaseName, key) {
			return value, nil
		}
	}

	return "", errors.New("No activity type found")

}
