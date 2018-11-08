package geo

import (
	"fmt"
	"time"
)

//MovingData contains moving data
type MovingData struct {
	Duration        float64
	Distance        float64
	MovingTime      float64
	StoppedTime     float64
	MovingDistance  float64
	StoppedDistance float64
	MaxSpeed        float64
	AverageSpeed    float64
	MaxPace         float64
	AveragePace     float64
	Points          []GPXPoint
	NumberValues    int
	SumAverageSpeed float64
	SumAveragePace  float64
}

func (md *MovingData) String(prefix string) string {
	t00, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.Duration)))
	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.MovingTime)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.StoppedTime)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.MaxPace*16.666666666667*60)))
	t04, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(md.AveragePace*16.666666666667*60)))
	var result string
	result = fmt.Sprintf("%s--- MovingData ---\n", prefix)
	result += fmt.Sprintf("%sDuration: %s\n", prefix, t00)
	result += fmt.Sprintf("%sDistance: %f km\n", prefix, md.Distance/1000.0)
	result += fmt.Sprintf("%sMoving Time: %s\n", prefix, t01)
	result += fmt.Sprintf("%sMoving Distance: %f km\n", prefix, md.MovingDistance/1000.0)
	result += fmt.Sprintf("%sStopped Time: %s\n", prefix, t02)
	result += fmt.Sprintf("%sStopped Distance: %f km\n", prefix, md.StoppedDistance/1000.0)
	result += fmt.Sprintf("%sMax Speed: %f m/sec -> %f km/h\n", prefix, md.MaxSpeed, md.MaxSpeed*3.6)
	result += fmt.Sprintf("%sAverage Speed: %f m/sec -> %f km/h\n", prefix, md.AverageSpeed, md.AverageSpeed*3.6)
	result += fmt.Sprintf("%sMax Pace: %f sec/m -> %s/km\n", prefix, md.MaxPace, t03)
	result += fmt.Sprintf("%sAverage Pace: %f sec/m -> %s/km\n", prefix, md.AveragePace, t04)
	result += fmt.Sprintf("%s------\n", prefix)
	return result
}

// SetMovingValues sets the mvoing and stopped values for MovingData (MovingTime, StoppedTime, MobingDistance, SzppedDistance, MaxSpeed, AverageSpeed)
func (md *MovingData) SetMovingValues(movingTime float64, stoppedTime float64, movingDistance float64, stoppedDistance float64, maxSpeed float64, averageSpeed float64, maxPace float64, averagePace float64, numberValues int) {
	md.MovingTime += movingTime
	md.StoppedTime += stoppedTime
	md.MovingDistance += movingDistance
	md.StoppedDistance += stoppedDistance

	md.Duration = md.MovingTime + md.StoppedTime
	md.Distance = md.MovingDistance + md.StoppedDistance

	if maxSpeed > md.MaxSpeed {
		md.MaxSpeed = maxSpeed
	}
	md.SumAverageSpeed += averageSpeed
	if maxPace > md.MaxPace {
		md.MaxPace = maxPace
	}
	md.SumAveragePace += averagePace
	md.NumberValues = numberValues
	md.AverageSpeed = md.SumAverageSpeed / float64(md.NumberValues)
	md.AveragePace = md.SumAveragePace / float64(md.NumberValues)

}
