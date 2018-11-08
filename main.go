package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpx"
)

func readFiles() {
	defer timeTrack(time.Now(), "readFiles")
	vincenty := geo.Vincenty{
		SigmaMultiplier: 1.644854,
		OneDegree:       1000.0 * 10000.8 / 90.0,
		EarthRadius:     6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:      1 / 298.257223563,
		SemiMinorAxisB:  6356752.314245,
		Epsilon:         1e-12,
		MaxIterations:   200,
	}

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
		distance        float64
		duration        float64
		movingDistance  float64
		movingTime      float64
		stoppedDistance float64
		stoppedTime     float64
	)
	for _, file := range files {
		if file.IsDir() == false {
			gpxDoc, err := gpx.ParseFile(filepath.Join(currentDirectory, "test", "gpx_all", file.Name()), &vincenty)
			if err != nil {
				panic(err)
			}
			md := gpxDoc.MovingData
			distance += md.Distance
			duration += md.Duration
			movingDistance += md.MovingDistance
			movingTime += md.MovingTime
			stoppedDistance += md.StoppedDistance
			stoppedTime += md.StoppedTime
		}
	}

	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(duration)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(movingTime)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(stoppedTime)))

	var result string
	result = fmt.Sprintf("--- GPX Data ---\n")
	result += fmt.Sprintf("Distance: %f\n", distance/1000.0)
	result += fmt.Sprintf("Duration Time: %s\n", t01)
	result += fmt.Sprintf("Moving Distance: %f km\n", movingDistance/1000.0)
	result += fmt.Sprintf("Moving Time Time: %s\n", t02)
	result += fmt.Sprintf("Stopped Distance: %f km\n", stoppedDistance/1000.0)
	result += fmt.Sprintf("Stopped Time Time: %s\n", t03)
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
