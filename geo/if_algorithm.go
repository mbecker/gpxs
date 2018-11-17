package geo

// Algorithm interface defines the customs funcs
type Algorithm interface {
	// Checks the activity tpe by the name (gpx, track, segment) and returns the activity tpye defined by an int
	CheckActivityType(lowerCaseName string) (string, error)

	// Should the normalization methof of standard deviation be used to determine which points belongs to moving and which one to stopped time/distance
	ShouldStandardDeviation() bool
	// Sigma defines the multiplier for the standard deviation to define x1 and x2 in which all points should be to define moving time/distance
	Sigma() float64

	// Return statement: Sshould the Point be included in the "MovingTime, MovingDistance" or "StoppedTime, StoppedDistance"
	// The gxPoint.Point.MovingData must be set in this func (!)
	CustomMovingPoints(gpxPoint *GPXPoint, previousGPXPoint *GPXPoint, algorithm Algorithm) error

	// Returns the duration between previous and actual point
	Duration(p1 Point, previousPoint Point) (float64, error)
	// Returns the duration between previous and actual point
	Distance(p1 *Point, previousPoint *Point) (float64, error)
	// Returns the speed between previous and actual point
	Speed(distance float64, duration float64) (float64, error)
	// Returns the pace between previous and actual point
	Pace(distance float64, duration float64) (float64, error)
}
