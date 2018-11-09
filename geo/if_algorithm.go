package geo

// Algorithm interface defines the funcs used to define the duration, distance, ... between two points
type Algorithm interface {
	ShouldStandardDeviation() bool
	Sigma() float64 // Sigma defines the multiplier for the standard deviation to define x1 and x2 in which all points should be to define movingData
	Duration(p1 *Point, previousPoint *Point) (float64, error)
	Distance(p1 *Point, previousPoint *Point) (float64, error)
	Speed(distance float64, duration float64) (float64, error)
	Pace(distance float64, duration float64) (float64, error)
	CustomMovingPoints(gpxPoint *GPXPoint, previousGPXPoint *GPXPoint, algorithm Algorithm) error // return statement: (should the Point be included in the "MovingTime, MovingDistance" or "StoopedTime, StoppedDistance", MovingData to be used for "MovingTime, MovingDistance", Duration of the point used for StoppedDuration, Distance of the point used for StoppedDistance
}
