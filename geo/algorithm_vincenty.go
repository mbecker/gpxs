package geo

import (
	"errors"
	"fmt"
	"math"
)

// Vincenty implements the Vincenty formula to calculate the distance
type Vincenty struct {
	ShouldStandardDeviationBeUsed bool    // Should the standard deviation be used to determine which points are used for calculation
	SigmaMultiplier               float64 // Define the sima standard deviation
	OneDegree                     float64
	EarthRadius                   float64
	Flattening                    float64
	SemiMinorAxisB                float64
	Epsilon                       float64
	MaxIterations                 int
}

// ShouldStandardDeviation (Vincenty) returns if the standard deviation should be used or not
func (v *Vincenty) ShouldStandardDeviation() bool {
	return v.ShouldStandardDeviationBeUsed
}

// Sigma (Vincenty) returns the sigma for the standard deviation
func (v *Vincenty) Sigma() float64 {
	return v.SigmaMultiplier
}

// Duration (Vincenty) returns the time.Duration from point p1 to previousPoint in sec
func (v *Vincenty) Duration(p1 *Point, previousPoint *Point) (float64, error) {
	return p1.Timestamp.Sub(previousPoint.Timestamp).Seconds(), nil
}

// CustomMovingPoints (Vincenty) defines which points should be used for "Moving"Time/Distance and if the it's set the new gpxPoint.Point Data
func (v *Vincenty) CustomMovingPoints(gpxPoint *GPXPoint, previousGPXPoint *GPXPoint, algorithm Algorithm) error {

	/* 	Define which points should be used; if a point should be used for calculation then set it's new values like Duration, Distance, Speed, etc.
	Here we use the set the new value for the points which used for "Moving"Time/Distanc
	*/

	// speed < 1 m/s
	if gpxPoint.Speed < 1.0 {
		fmt.Println("???")
		return errors.New("Point Speed below threshold")
	}
	gpxPoint.Point.SetPointData(previousGPXPoint.Point, algorithm)
	return nil
}

// Vincenty formula
func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Distance (Vincenty) returns the geographical distance in km between the points p1 (lat1, long1) and p2 (lat2, long2) using Vincenty's inverse formula.
// The surface of the Earth is approximated by the WGS-84 ellipsoid.
// This method may fail to converge for nearly antipodal points.
// https://github.com/asmarques/geodist/blob/master/vincenty.go
func (v *Vincenty) Distance(p1 *Point, previousPoint *Point) (float64, error) {
	if p1.Latitude == previousPoint.Latitude && p1.Longitude == previousPoint.Longitude {
		return 0, nil
	}

	U1 := math.Atan((1 - v.Flattening) * math.Tan(toRadians(p1.Latitude)))
	U2 := math.Atan((1 - v.Flattening) * math.Tan(toRadians(previousPoint.Latitude)))
	L := toRadians(previousPoint.Longitude - p1.Longitude)
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)
	lambda := L

	result := math.NaN()

	for i := 0; i < v.MaxIterations; i++ {
		curLambda := lambda
		sinSigma := math.Sqrt(math.Pow(cosU2*math.Sin(lambda), 2) +
			math.Pow(cosU1*sinU2-sinU1*cosU2*math.Cos(lambda), 2))
		cosSigma := sinU1*sinU2 + cosU1*cosU2*math.Cos(lambda)
		sigma := math.Atan2(sinSigma, cosSigma)
		sinAlpha := (cosU1 * cosU2 * math.Sin(lambda)) / math.Sin(sigma)
		cosSqrAlpha := 1 - math.Pow(sinAlpha, 2)
		cos2sigmam := 0.0
		if cosSqrAlpha != 0 {
			cos2sigmam = math.Cos(sigma) - ((2 * sinU1 * sinU2) / cosSqrAlpha)
		}
		C := (v.Flattening / 16) * cosSqrAlpha * (4 + v.Flattening*(4-3*cosSqrAlpha))
		lambda = L + (1-C)*v.Flattening*sinAlpha*(sigma+C*sinSigma*(cos2sigmam+C*cosSigma*(-1+2*math.Pow(cos2sigmam, 2))))

		if math.Abs(lambda-curLambda) < v.Epsilon {
			uSqr := cosSqrAlpha * ((math.Pow(v.EarthRadius, 2) - math.Pow(v.SemiMinorAxisB, 2)) / math.Pow(v.SemiMinorAxisB, 2))
			k1 := (math.Sqrt(1+uSqr) - 1) / (math.Sqrt(1+uSqr) + 1)
			A := (1 + (math.Pow(k1, 2) / 4)) / (1 - k1)
			B := k1 * (1 - (3*math.Pow(k1, 2))/8)

			deltaSigma := B * sinSigma * (cos2sigmam + (B/4)*(cosSigma*(-1+2*math.Pow(cos2sigmam, 2))-
				(B/6)*cos2sigmam*(-3+4*math.Pow(sinSigma, 2))*(-3+4*math.Pow(cos2sigmam, 2))))
			s := v.SemiMinorAxisB * A * (sigma - deltaSigma)
			result = s

			break
		}
	}

	if math.IsNaN(result) {
		return result, fmt.Errorf("Failed to converge for Point(Latitude: %f, Longitude: %f) and Point(Latitude: %f, Longitude: %f)", p1.Latitude, p1.Longitude, previousPoint.Latitude, previousPoint.Latitude)
	}

	return result, nil
}

// Speed (Vincenty) returns the speed in m/s
func (v *Vincenty) Speed(distance float64, duration float64) (float64, error) {
	return distance / duration, nil
}

// Pace (Vincenty) returns the pace in s/m
func (v *Vincenty) Pace(distance float64, duration float64) (float64, error) {
	return duration / distance, nil
}
