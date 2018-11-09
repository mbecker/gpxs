package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
	"github.com/olekukonko/tablewriter"
)

type AlgStruct struct {
	Name string
	Alg  geo.Algorithm
}

func parseFiles(fileDirectory string, files []os.FileInfo, alg geo.Algorithm, tableData [][]string, tableFooter [][]string) {
	start := time.Now()
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
			gpxDoc, err := gpxs.ParseFile(filepath.Join(fileDirectory, file.Name()), alg)
			if err != nil {
				panic(err)
			}

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
	elapsed := time.Since(start)

	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(duration)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(movingTime)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(stoppedTime)))

	t04, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackDuration)))
	t05, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackMovingTime)))
	t06, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackStoppedTime)))

	t07, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentDuration)))
	t08, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentMovingTime)))
	t09, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentStoppedTime)))

	tableData[0] = append(tableData[0], fmt.Sprintf("%d", countFiles)) // Metadata: # of files
	tableData[1] = append(tableData[1], "------")
	tableData[2] = append(tableData[2], fmt.Sprintf("%s", t01))                    // GPX Duration
	tableData[3] = append(tableData[3], fmt.Sprintf("%f", distance/1000.0))        // GPX Distance
	tableData[4] = append(tableData[4], fmt.Sprintf("%s", t02))                    // GPX Moving Time)
	tableData[5] = append(tableData[5], fmt.Sprintf("%s", t03))                    // GPX Stopped Time
	tableData[6] = append(tableData[6], fmt.Sprintf("%f", movingDistance/1000.0))  // GPX Moving Distance
	tableData[7] = append(tableData[7], fmt.Sprintf("%f", stoppedDistance/1000.0)) // GPX Stopped Distance
	tableData[8] = append(tableData[8], "------")
	tableData[9] = append(tableData[9], fmt.Sprintf("%s", t04))                           // Track Duration
	tableData[10] = append(tableData[10], fmt.Sprintf("%f", trackDistance/1000.0))        // Track Moving Distance
	tableData[11] = append(tableData[11], fmt.Sprintf("%s", t05))                         // Track Moving Time
	tableData[12] = append(tableData[12], fmt.Sprintf("%s", t06))                         // Track Stopped Time
	tableData[13] = append(tableData[13], fmt.Sprintf("%f", trackMovingDistance/1000.0))  // Track Moving Distance
	tableData[14] = append(tableData[14], fmt.Sprintf("%f", trackStoppedDistance/1000.0)) // Track Stopped Distance
	tableData[15] = append(tableData[15], "------")
	tableData[16] = append(tableData[16], fmt.Sprintf("%s", t07))                           // Segment Duration
	tableData[17] = append(tableData[17], fmt.Sprintf("%f", segmentDistance/1000.0))        // Segment Distance
	tableData[18] = append(tableData[18], fmt.Sprintf("%s", t08))                           // Segment Moving Time
	tableData[19] = append(tableData[19], fmt.Sprintf("%s", t09))                           // Segment  Stopped Time
	tableData[20] = append(tableData[20], fmt.Sprintf("%f", segmentMovingDistance/1000.0))  // Segment Moving Distance
	tableData[21] = append(tableData[21], fmt.Sprintf("%f", segmentStoppedDistance/1000.0)) // Segment  Stopped Distance

	tableFooter[0] = append(tableFooter[0], fmt.Sprintf("%s", elapsed))
}

func main() {

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

	vincentyWithoutStandardDeviation := geo.Vincenty{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               1.644854, // ~95%
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:                    1 / 298.257223563,
		SemiMinorAxisB:                6356752.314245,
		Epsilon:                       1e-12,
		MaxIterations:                 200,
	}

	algorithmStandardLength2d := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               1.644854, // ~95%
		ShouldHaversine:               true,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      false,
	}

	algorithmStandarLength2ddWithStandardDeviation := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               1.644854, // ~95%
		ShouldHaversine:               true,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      false,
	}

	algorithmStandardLength3d := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               1.644854, // ~95%
		ShouldHaversine:               true,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      true,
	}

	algorithmStandarLength3ddWithStandardDeviation := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               1.644854, // ~95%
		ShouldHaversine:               true,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      true,
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileDirectory := filepath.Join(currentDirectory, "test", "gpx_all")
	files, err := ioutil.ReadDir(fileDirectory)
	if err != nil {
		panic(err)
	}

	algorithms := []AlgStruct{
		AlgStruct{
			Name: "Vincenty w/o Standard Deviation",
			Alg:  &vincentyWithoutStandardDeviation,
		},
		AlgStruct{
			Name: "Vincenty with Standard Deviation",
			Alg:  &vincenty,
		},
		AlgStruct{
			Name: "Standard (length2D) W/o standard deviation",
			Alg:  &algorithmStandardLength2d,
		},
		AlgStruct{
			Name: "Standard (length2D) With standard deviation",
			Alg:  &algorithmStandarLength2ddWithStandardDeviation,
		},
		AlgStruct{
			Name: "Standard (length3D) W/o standard deviation",
			Alg:  &algorithmStandardLength3d,
		},
		AlgStruct{
			Name: "Standard (length3D) With standard deviation",
			Alg:  &algorithmStandarLength3ddWithStandardDeviation,
		},
	}

	tableFooter := [][]string{
		[]string{
			"Execution time",
		},
	}

	tableData := [][]string{
		[]string{
			"# of files",
		},
		[]string{
			"------",
		},
		[]string{
			"GPX Duration",
		},
		[]string{
			"GPX Distance",
		},
		[]string{
			"GPX Moving Time",
		},
		[]string{
			"GPX Stopped Time",
		},
		[]string{
			"GPX Moving Distance",
		},
		[]string{
			"GPX Stopped Distance",
		},
		[]string{
			"------",
		},
		[]string{
			"Track Duration",
		},
		[]string{
			"Track Distance",
		},
		[]string{
			"Track Moving Time",
		},
		[]string{
			"Track Stopped Time",
		},
		[]string{
			"Track Moving Distance",
		},
		[]string{
			"Track Stopped Distance",
		},
		[]string{
			"------",
		},
		[]string{
			"Segment Duration",
		},
		[]string{
			"Segment Distance",
		},
		[]string{
			"Segment Moving Time",
		},
		[]string{
			"Segment Stopped Time",
		},
		[]string{
			"Segment Moving Distance",
		},
		[]string{
			"Segment Stopped Distance",
		},
	}

	for _, alg := range algorithms {
		parseFiles(fileDirectory, files, alg.Alg, tableData, tableFooter)
	}

	// readFiles()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetCaption(true, "gpxs Benchmark - The execution with the Vincenty forula takes longer due to more caluclation steps")

	var header []string
	header = make([]string, 0)
	header = append(header, "Type")
	for _, alg := range algorithms {
		header = append(header, alg.Name)
	}
	table.SetHeader(header)
	table.SetFooter(tableFooter[0])
	table.AppendBulk(tableData)
	table.Render() // Send output
}
