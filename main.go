package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
	"github.com/olekukonko/tablewriter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type AlgStruct struct {
	Name string
	Alg  geo.Algorithm
}

func executionTime(start time.Time, name string) {
	duration := time.Since(start)
	fmt.Printf("Execution time: %s - %s", name, duration)
}

func parseFilesForPlot(fileDirectory string, files []os.FileInfo, alg geo.Algorithm) {
	for _, file := range files {

		if file.IsDir() == false && filepath.Ext(file.Name()) == ".gpx" {
			gpxDoc, err := gpxs.ParseFile(filepath.Join(fileDirectory, file.Name()), alg)
			if err != nil {
				panic(err)
			}
			var extension = filepath.Ext(file.Name())
			var name = file.Name()[0 : len(file.Name())-len(extension)]
			fileName := fmt.Sprintf("%s-graph-%s.png", name, alg.String())
			graphPath := filepath.Join(fileDirectory, fileName)
			createPlot(gpxDoc, alg.String(), graphPath)
		}
	}
}

func createPlot(gpx *geo.GPX, title string, filePath string) {

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Title.Text = title
	p.X.Label.Text = "Distance"
	p.Y.Label.Text = "Speed (m/s)"

	points := graphPoints(gpx)

	// Make a line plotter with points and set its style.
	lpLineOverallPoints, lpPointsOverallPoints, err := plotter.NewLinePoints(points[2])
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
		return
	}
	// lpLineOverallPoints.Color = color.RGBA{G: 255, A: 255}
	// lpPointsOverallPoints.Shape = draw.PyramidGlyph{}
	// lpPointsOverallPoints.Color = color.RGBA{R: 255, A: 255}

	// Make a scatter plotter and set its style.
	scatterStoppedPoints, err := plotter.NewScatter(points[3])
	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
		return
	}
	scatterStoppedPoints.GlyphStyle.Color = color.RGBA{242, 60, 80, 0xff}
	scatterStoppedPoints.Color = color.RGBA{242, 60, 80, 0xff}
	scatterStoppedPoints.GlyphStyle.Shape = draw.CircleGlyph{}
	scatterStoppedPoints.GlyphStyle.Radius = 4
	// Add the plotters to the plot, with a legend entry for each
	p.Add(lpLineOverallPoints, lpPointsOverallPoints, scatterStoppedPoints)
	p.Legend.Add("Overall Points", lpLineOverallPoints, lpPointsOverallPoints)
	p.Legend.Add("Stopped Points", scatterStoppedPoints)

	// Standard Deviation
	if gpx.MovementStats.SD.Valid {
		// Standard Deviation lower border
		lpLineX1, _, err := plotter.NewLinePoints(points[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		lpLineX1.Color = color.RGBA{242, 60, 80, 0xff}

		// Standard Deviation higher border
		lpLineX2, _, err := plotter.NewLinePoints(points[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		lpLineX2.Color = color.RGBA{242, 60, 80, 0xff}
		p.Add(lpLineX1, lpLineX2)
		p.Legend.Add("SD y1", lpLineX1)
		p.Legend.Add("SD y2", lpLineX2)
	}

	// Save the plot to a PNG file.
	if err := p.Save(21*vg.Inch, 12*vg.Inch, filePath); err != nil {
		fmt.Println(err)
		return
	}

}

// graphPoints returns the points for moving and stopped points, and adds the x1/x1 (here y1/y2) line of standard deviation
func graphPoints(gpx *geo.GPX) map[int]plotter.XYs {

	ptX1 := make(plotter.XYs, gpx.PointsCount)
	ptX2 := make(plotter.XYs, gpx.PointsCount)
	if gpx.MovementStats.SD.Valid {
		ptX1[0].X = 0
		ptX1[0].Y = gpx.MovementStats.SD.X1

		ptX2[0].X = 0
		ptX2[0].Y = gpx.MovementStats.SD.X2
	}

	ptAll := make(plotter.XYs, gpx.PointsCount)
	ptStopped := make(plotter.XYs, gpx.PointsCount)

	var distance float64
	i := 0
	for _, track := range gpx.Tracks {
		for _, segment := range track.Segments {
			for x, point := range segment.Points {
				distance += (point.Distance / 1000)

				// Standard Deviation
				if gpx.MovementStats.SD.Valid {
					ptX1[i].X = distance
					ptX1[i].Y = gpx.MovementStats.SD.X1

					ptX2[i].X = distance
					ptX2[i].Y = gpx.MovementStats.SD.X2
				}

				ptAll[i].X = distance
				ptAll[i].Y = point.Distance / point.Duration
				if !point.IsMoving {
					ptStopped[i].X = distance
					ptStopped[i].Y = point.Distance / point.Duration
				}
				if x == 0 {
					ptAll[0].X = 0
					ptAll[0].Y = segment.Points[x+1].Distance / segment.Points[x+1].Duration // Uhh Dangerous
				}
				i++
			}
		}
	}

	m := make(map[int]plotter.XYs)
	m[0] = ptX1
	m[1] = ptX2
	m[2] = ptAll
	m[3] = ptStopped

	return m
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
		segmentMaxPace         float64
		segmentAveragePace     float64
		segmentMaxSpeed        float64
		segmentAverageSpeed    float64

		gpxDocs []*geo.GPX
	)

	for _, file := range files {

		if file.IsDir() == false && filepath.Ext(file.Name()) == ".gpx" {
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
					segmentMaxPace = segment.MovementStats.StoppedData.MaxPace
					segmentAveragePace = segment.MovementStats.StoppedData.AveragePace
					segmentMaxSpeed = segment.MovementStats.StoppedData.MaxSpeed
					segmentAverageSpeed = segment.MovementStats.StoppedData.AverageSpeed
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
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentMaxSpeed)) // Segment max pace
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentAverageSpeed)) // Segment max pace
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentMaxPace)) // Segment max pace
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%v", segmentAveragePace)) // Segment max pace
	x++
	tableData[x] = append(tableData[x], "------")
	x++
	tableData[x] = append(tableData[x], fmt.Sprintf("%s", elapsed))

}

func main() {

	// customAlgorithm := CustomAlgorithm{
	// 	CustomParameter: 100.9,
	// }

	var sigmaMultiplier = 1.959964

	vincentyWithoutStandardDeviation := geo.Vincenty{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               sigmaMultiplier, // ~90%
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:                    1 / 298.257223563,
		SemiMinorAxisB:                6356752.314245,
		Epsilon:                       1e-12,
		MaxIterations:                 200,
		Name:                          "VincentySpeedThreshold",
	}

	vincenty := geo.Vincenty{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               sigmaMultiplier, // ~95%
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137, // WGS-84 ellipsoid; See https://en.wikipedia.org/wiki/World_Geodetic_System
		Flattening:                    1 / 298.257223563,
		SemiMinorAxisB:                6356752.314245,
		Epsilon:                       1e-12,
		MaxIterations:                 200,
		Name:                          "VincentySD",
	}

	algorithmStandardLength2d := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               sigmaMultiplier, // ~95%
		ShouldHaversine:               false,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      false,
		Name:                          "gpxgoLength2d",
	}

	algorithmStandarLength2ddWithStandardDeviation := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               sigmaMultiplier, // ~95%
		ShouldHaversine:               false,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      false,
		Name:                          "gpxgoLength2dSD",
	}

	algorithmStandardLength3d := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: false,
		SigmaMultiplier:               sigmaMultiplier, // ~95%
		ShouldHaversine:               false,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      true,
		Name:                          "gpxgoLength3d",
	}

	algorithmStandarLength3ddWithStandardDeviation := geo.AlgorithmStandard{
		ShouldStandardDeviationBeUsed: true,
		SigmaMultiplier:               sigmaMultiplier, // ~95%
		ShouldHaversine:               false,
		OneDegree:                     1000.0 * 10000.8 / 90.0,
		EarthRadius:                   6378137,
		Should3D:                      true,
		Name:                          "gpxgoLength3dSD",
	}

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
			"Segment Max Speed",
		},
		[]string{
			"Segment Average Speed",
		},
		[]string{
			"Segment Max Pace",
		},
		[]string{
			"Segment Average Pace",
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
		parseFilesForPlot(fileDirectory, files, alg.Alg)
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
