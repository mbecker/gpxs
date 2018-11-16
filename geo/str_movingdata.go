package geo

import (
	"database/sql/driver"
	"fmt"
	"math"
	"time"
)

// NullTime from https://github.com/lib/pq/blob/8c6ee72f3e6bcb1542298dd5f76cb74af9742cec/encode.go#L586
type NullTime struct {
	Time  *time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTime) SetTime(time *time.Time) {
	if !time.IsZero() {
		nt.Time = time
		nt.Valid = true
	}
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(*time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MovementStats contains moving data
type MovementStats struct {
	OverallData MovementData
	MovingData  MovementData
	StoppedData MovementData
}

// MovementData represent the data/stats of 'overall', 'moving' and 'stopped'
type MovementData struct {
	Count int // The count of points in this structs

	StartTime NullTime
	EndTime   NullTime

	Duration float64 // The duration of the gpx, track, segment in sec
	Distance float64 // The distance of the gpx, track, segment defined by Alogorithm

	MaxSpeed     float64
	AverageSpeed float64

	MaxPace     float64
	AveragePace float64

	MinLatitude  float64
	MaxLatitude  float64
	MinLongitude float64
	MaxLongitude float64
	MinEvelation float64
	MaxEvelation float64

	Points []*GPXPoint
}

func (ms *MovementStats) String() string {
	result := fmt.Sprintf("--- Movemenet Stats ---")
	result += ms.OverallData.String("OverallData", " ")
	result += ms.MovingData.String("OverallData", " ")
	result += ms.StoppedData.String("OverallData", " ")
	return result
}

func (md *MovementData) String(title string, prefix string) string {

	t00, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.Duration)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.MaxPace*16.666666666667*60)))
	t04, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.AveragePace*16.666666666667*60)))
	var result string
	result = fmt.Sprintf("%s--- %s ---\n", prefix, title)
	result += fmt.Sprintf("%sStartTime: %s\n", md.StartTime)
	result += fmt.Sprintf("%EndTime: %s\n", md.EndTime)
	result += fmt.Sprintf("%sDuration: %s\n", prefix, t00)
	result += fmt.Sprintf("%sDistance: %f km\n", prefix, md.Distance/1000.0)
	result += fmt.Sprintf("%sMax Speed: %f m/sec -> %f km/h\n", prefix, md.MaxSpeed, md.MaxSpeed*3.6)
	result += fmt.Sprintf("%sAverage Speed: %f m/sec -> %f km/h\n", prefix, md.AverageSpeed, md.AverageSpeed*3.6)
	result += fmt.Sprintf("%sMax Pace: %f sec/m -> %s/km\n", prefix, md.MaxPace, t03)
	result += fmt.Sprintf("%sAverage Pace: %f sec/m -> %s/km\n", prefix, md.AveragePace, t04)
	result += fmt.Sprintf("%MinLatitude: %f\n", md.MinLatitude)
	result += fmt.Sprintf("%MaxLatitude: %f\n", md.MaxLatitude)
	result += fmt.Sprintf("%MinLongitude: %f\n", md.MinLongitude)
	result += fmt.Sprintf("%MaxLongitude: %f\n", md.MaxLongitude)
	result += fmt.Sprintf("%MinEvelation: %f\n", md.MinEvelation)
	result += fmt.Sprintf("%MaxEvelation: %f\n", md.MaxEvelation)
	result += fmt.Sprintf("%s------\n", prefix)
	return result
}

// SetMovingValues sets the mvoing and stopped values for MovingData (MovingTime, StoppedTime, MobingDistance, SzppedDistance, MaxSpeed, AverageSpeed)
// func (md *MovingData) SetMovingValues(movingTime float64, stoppedTime float64, movingDistance float64, stoppedDistance float64, maxSpeed float64, averageSpeed float64, maxPace float64, averagePace float64, numberValues int) {
// 	md.MovingTime += movingTime
// 	md.StoppedTime += stoppedTime
// 	md.MovingDistance += movingDistance
// 	md.StoppedDistance += stoppedDistance

// 	md.Duration = md.MovingTime + md.StoppedTime
// 	md.Distance = md.MovingDistance + md.StoppedDistance

// 	if maxSpeed > md.MaxSpeed {
// 		md.MaxSpeed = maxSpeed
// 	}
// 	md.SumAverageSpeed += averageSpeed
// 	if maxPace > md.MaxPace {
// 		md.MaxPace = maxPace
// 	}
// 	md.SumAveragePace += averagePace
// 	md.NumberValues = numberValues
// 	md.AverageSpeed = md.SumAverageSpeed / float64(md.NumberValues)
// 	md.AveragePace = md.SumAveragePace / float64(md.NumberValues)

// }

// SetValues set the values for MovementData
func (md *MovementData) SetValues(gpxPoint *GPXPoint, previousGpxPoint *GPXPoint, count int, alg Algorithm) {
	md.Count = count

	md.EndTime.SetTime(gpxPoint.Timestamp.Time)

	md.Duration += gpxPoint.Duration
	md.Distance += gpxPoint.Distance

	// Speed
	if gpxPoint.Speed > md.MaxSpeed {
		md.MaxSpeed = gpxPoint.Speed
	}
	if averageSpeed, err := alg.Speed(md.Distance, md.Duration); err == nil {
		md.AverageSpeed = averageSpeed
	}

	// Pace
	if md.MaxPace == 0 || (gpxPoint.Pace > 0.0 && gpxPoint.Pace < md.MaxPace) {
		md.MaxPace = gpxPoint.Pace
	}
	if averagePace, err := alg.Pace(md.Distance, md.Duration); err == nil {
		md.AveragePace = averagePace
	}

	// Latitude
	if (md.MinLatitude == 0 && math.Abs(gpxPoint.Latitude) > 0) || (math.Abs(md.MinLatitude) > math.Abs(gpxPoint.Latitude)) {
		md.MinLatitude = gpxPoint.Latitude
	}
	if (md.MaxLatitude == 0 && gpxPoint.Latitude > 0) || (md.MaxLatitude > gpxPoint.Latitude) {
		md.MaxLatitude = gpxPoint.Latitude
	}

	// Longitude
	if (md.MinLongitude == 0 && math.Abs(gpxPoint.Longitude) > 0) || (math.Abs(md.MinLongitude) > math.Abs(gpxPoint.Longitude)) {
		md.MinLongitude = gpxPoint.Longitude
	}
	if (md.MaxLongitude == 0 && gpxPoint.Longitude > 0) || (md.MaxLongitude > gpxPoint.Longitude) {
		md.MaxLongitude = gpxPoint.Longitude
	}

	// Evelation
	if (md.MinEvelation == 0 && math.Abs(gpxPoint.Elevation.Value()) > 0) || (math.Abs(md.MinEvelation) > math.Abs(gpxPoint.Elevation.Value())) {
		md.MinEvelation = gpxPoint.Elevation.Value()
	}
	if (md.MaxEvelation == 0 && gpxPoint.Elevation.Value() > 0) || (md.MaxEvelation > gpxPoint.Elevation.Value()) {
		md.MaxEvelation = gpxPoint.Elevation.Value()
	}

	md.Points = append(md.Points, gpxPoint)
}

func (md *MovementData) SetValuesFromMovementData(movementData *MovementData, count int, alg Algorithm) {
	md.Count = count

	// Before reports whether the time instant md.EndTime.Time is before u.
	if md.EndTime.Valid && movementData.EndTime.Valid && md.EndTime.Time.Before(*movementData.EndTime.Time) {
		md.EndTime = movementData.EndTime
	}

	md.Duration += movementData.Duration
	md.Distance += movementData.Distance

	// Speed
	if movementData.MaxSpeed > md.MaxSpeed {
		md.MaxSpeed = movementData.MaxSpeed
	}
	if averageSpeed, err := alg.Speed(md.Distance, md.Duration); err == nil {
		md.AverageSpeed = averageSpeed
	}

	// Pace
	if md.MaxPace == 0 || (movementData.MaxPace > 0.0 && movementData.MaxPace < md.MaxPace) {
		md.MaxPace = movementData.MaxPace
	}
	if averagePace, err := alg.Pace(md.Distance, md.Duration); err == nil {
		md.AveragePace = averagePace
	}

	// Latitude
	if (md.MinLatitude == 0 && math.Abs(movementData.MinLatitude) > 0) || (math.Abs(md.MinLatitude) > math.Abs(movementData.MinLatitude)) {
		md.MinLatitude = movementData.MinLatitude
	}
	if (md.MaxLatitude == 0 && movementData.MaxLatitude > 0) || (md.MaxLatitude > movementData.MaxLatitude) {
		md.MaxLatitude = movementData.MaxLatitude
	}

	// Longitude
	if (md.MinLongitude == 0 && math.Abs(movementData.MinLongitude) > 0) || (math.Abs(md.MinLongitude) > math.Abs(movementData.MinLongitude)) {
		md.MinLongitude = movementData.MinLongitude
	}
	if (md.MaxLongitude == 0 && movementData.MaxLongitude > 0) || (md.MaxLongitude > movementData.MaxLongitude) {
		md.MaxLongitude = movementData.MaxLongitude
	}

	// Evelation
	if (md.MinEvelation == 0 && math.Abs(movementData.MinEvelation) > 0) || (math.Abs(md.MinEvelation) > math.Abs(movementData.MinEvelation)) {
		md.MinEvelation = movementData.MinEvelation
	}
	if (md.MaxEvelation == 0 && movementData.MaxEvelation > 0) || (md.MaxEvelation > movementData.MaxEvelation) {
		md.MaxEvelation = movementData.MaxEvelation
	}

	for _, point := range movementData.Points {
		md.Points = append(md.Points, point)
	}
}
