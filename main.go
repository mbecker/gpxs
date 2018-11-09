package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
)

func readFiles() {
	defer timeTrack(time.Now(), "readFiles")
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

	// algorithmStandard := geo.AlgorithmStandard{
	// 	ShouldStandardDeviationBeUsed: true,
	// 	SigmaMultiplier:               1.644854, // ~95%
	// 	ShouldHaversine:               true,
	// 	OneDegree:                     1000.0 * 10000.8 / 90.0,
	// 	EarthRadius:                   6378137,
	// 	Should3D:                      false,
	// }

	/** Benchmark **/
	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(filepath.Join(currentDirectory, "test", "gpx_all"))
	if err != nil {
		panic(err)
	}
	var (
		countFiles      int
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
	for _, file := range files {
		if file.IsDir() == false {
			gpxDoc, err := gpxs.ParseFile(filepath.Join(currentDirectory, "test", "gpx_all", file.Name()), &vincenty)
			if err != nil {
				panic(err)
			}
			// c := make(chan gpx.ChannelStruct)
			// go gpx.ParseFileChannel(filepath.Join(currentDirectory, "test", "gpx_all", file.Name()), &vincenty, c)
			// gpxChannelReturn := <-c
			// if gpxChannelReturn.Error != nil {
			// 	panic(gpxChannelReturn.Error)
			// }
			// gpxDoc := gpxChannelReturn.GPX

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

			countFiles++
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

	var result string
	result = fmt.Sprintf("--- GPX Files ---\n")
	result += fmt.Sprintf("# of Files: %d\n", countFiles)
	result += fmt.Sprintf("--- GPX ---\n")
	result += fmt.Sprintf("Distance: %f\n", distance/1000.0)
	result += fmt.Sprintf("Duration Time: %s\n", t01)
	result += fmt.Sprintf("Moving Distance: %f km\n", movingDistance/1000.0)
	result += fmt.Sprintf("Moving Time Time: %s\n", t02)
	result += fmt.Sprintf("Stopped Distance: %f km\n", stoppedDistance/1000.0)
	result += fmt.Sprintf("Stopped Time Time: %s\n", t03)

	result += fmt.Sprintf("--- Tracks ---\n")
	result += fmt.Sprintf("Distance: %f\n", trackDistance/1000.0)
	result += fmt.Sprintf("Duration Time: %s\n", t04)
	result += fmt.Sprintf("Moving Distance: %f km\n", trackMovingDistance/1000.0)
	result += fmt.Sprintf("Moving Time Time: %s\n", t05)
	result += fmt.Sprintf("Stopped Distance: %f km\n", trackStoppedDistance/1000.0)
	result += fmt.Sprintf("Stopped Time Time: %s\n", t06)

	result += fmt.Sprintf("--- Segment ---\n")
	result += fmt.Sprintf("Distance: %f\n", segmentDistance/1000.0)
	result += fmt.Sprintf("Duration Time: %s\n", t07)
	result += fmt.Sprintf("Moving Distance: %f km\n", segmentMovingDistance/1000.0)
	result += fmt.Sprintf("Moving Time Time: %s\n", t08)
	result += fmt.Sprintf("Stopped Distance: %f km\n", segmentStoppedDistance/1000.0)
	result += fmt.Sprintf("Stopped Time Time: %s\n", t09)

	result += fmt.Sprintf("------\n")
	fmt.Print(result)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	readFiles()
}
