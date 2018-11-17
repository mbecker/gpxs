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

func parseFiles(fileDirectory string, files []os.FileInfo, alg geo.Algorithm, tableData [][]string) {
	start := time.Now()
	var (
		countFiles      int
		distance        float64
		duration        float64
		movingDistance  float64
		movingTime      float64
		stoppedDistance float64
		stoppedTime     float64
		startTime       *time.Time
		endTime         *time.Time

		trackDistance        float64
		trackDuration        float64
		trackMovingDistance  float64
		trackMovingTime      float64
		trackStoppedDistance float64
		trackStoppedTime     float64
		trackStartTime       *time.Time
		trackEndTime         *time.Time

		segmentDistance        float64
		segmentDuration        float64
		segmentMovingDistance  float64
		segmentMovingTime      float64
		segmentStoppedDistance float64
		segmentStoppedTime     float64
		segmentStartTime       *time.Time
		segmentEndTime         *time.Time

		gpxDocs []*geo.GPX
	)

	for _, file := range files {

		if file.IsDir() == false {
			gpxDoc, err := gpxs.ParseFile(filepath.Join(fileDirectory, file.Name()), alg)
			if err != nil {
				panic(err)
			}
			gpxDocs = append(gpxDocs, gpxDoc)

			// GPX
			md := gpxDoc.MovementStats.OverallData
			distance += md.Distance
			duration += md.Duration
			movingDistance += gpxDoc.MovementStats.MovingData.Distance
			movingTime += gpxDoc.MovementStats.MovingData.Duration
			stoppedDistance += gpxDoc.MovementStats.StoppedData.Distance
			stoppedTime += gpxDoc.MovementStats.StoppedData.Duration
			startTime = gpxDoc.MovementStats.OverallData.StartTime.Time
			endTime = gpxDoc.MovementStats.OverallData.EndTime.Time

			// Tracks
			for _, track := range gpxDoc.Tracks {
				trackMd := track.MovementStats.OverallData
				trackDistance += trackMd.Distance
				trackDuration += trackMd.Duration
				trackMovingDistance += track.MovementStats.MovingData.Distance
				trackMovingTime += track.MovementStats.MovingData.Duration
				trackStoppedDistance += track.MovementStats.StoppedData.Distance
				trackStoppedTime += track.MovementStats.StoppedData.Duration
				trackStartTime = track.MovementStats.OverallData.StartTime.Time
				trackEndTime = track.MovementStats.OverallData.EndTime.Time

				// Segments
				for _, segment := range track.Segments {
					segmentMd := segment.MovementStats.OverallData
					segmentDistance += segmentMd.Distance
					segmentDuration += segmentMd.Duration
					segmentMovingDistance += segment.MovementStats.MovingData.Distance
					segmentMovingTime += segment.MovementStats.MovingData.Duration
					segmentStoppedDistance += segment.MovementStats.StoppedData.Distance
					segmentStoppedTime += segment.MovementStats.StoppedData.Duration
					segmentStartTime = segment.MovementStats.OverallData.StartTime.Time
					segmentEndTime = segment.MovementStats.OverallData.EndTime.Time
				}
			}

			countFiles++
		}
	}
	elapsed := time.Since(start)
	var gpxDocsTypeString string
	for _, doc := range gpxDocs {
		gpxDocsTypeString += fmt.Sprintf("%s, ", doc.Type)
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

	x := 0
	tableData[0] = append(tableData[0], fmt.Sprintf("%d", countFiles)) // Metadata: # of files
	x++
	tableData[x] = append(tableData[x], "------")
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", gpxDocsTypeString)) // GPX Duration
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t01)) // GPX Duration
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", distance/1000.0)) // GPX Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t02)) // GPX Moving Time)
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t03)) // GPX Stopped Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", movingDistance/1000.0)) // GPX Moving Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", stoppedDistance/1000.0)) // GPX Stopped Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", startTime)) // GPX Start Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", endTime)) // GPX Start Time
	x++
	tableData[x] = append(tableData[x], "------")
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t04)) // Track Duration
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", trackDistance/1000.0)) // Track Moving Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t05)) // Track Moving Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t06)) // Track Stopped Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", trackMovingDistance/1000.0)) // Track Moving Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", trackStoppedDistance/1000.0)) // Track Stopped Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", trackStartTime)) // Track Start Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", trackEndTime)) // Track Start Time
	x++
	tableData[x] = append(tableData[x], "------")
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t07)) // Segment Duration
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", segmentDistance/1000.0)) // Segment Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t08)) // Segment Moving Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", t09)) // Segment  Stopped Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", segmentMovingDistance/1000.0)) // Segment Moving Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%f", segmentStoppedDistance/1000.0)) // Segment  Stopped Distance
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentStartTime)) // Track Start Time
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentEndTime)) // Track Start Time
	x++
	tableData[x] = append(tableData[x], "------")
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", elapsed))

}

func main() {

	// customAlgorithm := CustomAlgorithm{
	// 	CustomParameter: 100.9,
	// }

	vincentyWithoutStandardDeviation := geo.Vincenty{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               1.644854, // ~95%
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:                    1 / 298.257223563,
		SemiMinorAxisB:                6356752.314245,
		Epsilon:                       1e-12,
		MaxIterations:                 200,
	}

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

	// algorithmStandardLength2d := geo.AlgorithmStandard{
	// 	ShouldStandardDeviationBeUsed: false,
	// 	SigmaMultiplier:               1.644854, // ~95%
	// 	ShouldHaversine:               false,
	// 	OneDegree:                     1000.0 * 10000.8 / 90.0,
	// 	EarthRadius:                   6378137,
	// 	Should3D:                      false,
	// }

	// algorithmStandarLength2ddWithStandardDeviation := geo.AlgorithmStandard{
	// 	ShouldStandardDeviationBeUsed: true,
	// 	SigmaMultiplier:               1.644854, // ~95%
	// 	ShouldHaversine:               false,
	// 	OneDegree:                     1000.0 * 10000.8 / 90.0,
	// 	EarthRadius:                   6378137,
	// 	Should3D:                      false,
	// }

	// algorithmStandardLength3d := geo.AlgorithmStandard{
	// 	ShouldStandardDeviationBeUsed: false,
	// 	SigmaMultiplier:               1.644854, // ~95%
	// 	ShouldHaversine:               false,
	// 	OneDegree:                     1000.0 * 10000.8 / 90.0,
	// 	EarthRadius:                   6378137,
	// 	Should3D:                      true,
	// }

	// algorithmStandarLength3ddWithStandardDeviation := geo.AlgorithmStandard{
	// 	ShouldStandardDeviationBeUsed: true,
	// 	SigmaMultiplier:               1.644854, // ~95%
	// 	ShouldHaversine:               false,
	// 	OneDegree:                     1000.0 * 10000.8 / 90.0,
	// 	EarthRadius:                   6378137,
	// 	Should3D:                      true,
	// }

	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileDirectory := filepath.Join(currentDirectory, "test", "gpx_files")
	files, err := ioutil.ReadDir(fileDirectory)
	if err != nil {
		panic(err)
	}

	algorithms := []AlgStruct{
		// AlgStruct{
		// 	Name: "Custom",
		// 	Alg:  &customAlgorithm,
		// },
		AlgStruct{
			Name: "Vinc. Speed Threshold (w/o SD)",
			Alg:  &vincentyWithoutStandardDeviation,
		},
		AlgStruct{
			Name: "Vinc. SD",
			Alg:  &vincenty,
		},
		// AlgStruct{
		// 	Name: "Standard (length2D) W/o standard deviation",
		// 	Alg:  &algorithmStandardLength2d,
		// },
		// AlgStruct{
		// 	Name: "Standard (length2D) With standard deviation",
		// 	Alg:  &algorithmStandarLength2ddWithStandardDeviation,
		// },
		// AlgStruct{
		// 	Name: "Standard (length3D) W/o standard deviation",
		// 	Alg:  &algorithmStandardLength3d,
		// },
		// AlgStruct{
		// 	Name: "Standard (length3D) With standard deviation",
		// 	Alg:  &algorithmStandarLength3ddWithStandardDeviation,
		// },
	}

	tableData := [][]string{
		[]string{
			"# of files",
		},
		[]string{
			"------",
		},
		[]string{
			"GPX Type",
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
			"GPX Start Time",
		},
		[]string{
			"GPX End Time",
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
			"Track Start Time",
		},
		[]string{
			"Track End Time",
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
		[]string{
			"Segment Start Time",
		},
		[]string{
			"Segment End Time",
		},
		[]string{
			"------",
		},
		[]string{
			"Execution time",
		},
	}

	for _, alg := range algorithms {
		parseFiles(fileDirectory, files, alg.Alg, tableData)
	}

	// readFiles()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	var header []string
	header = make([]string, 0)
	header = append(header, "Type")
	for _, alg := range algorithms {
		header = append(header, alg.Name)
	}
	table.SetHeader(header)
	table.AppendBulk(tableData)
	table.Render() // Send output
}
