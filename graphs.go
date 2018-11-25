package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func parseFilesForPlot(gpxFilesDirectory string, graphFilesDirectory string, files []os.FileInfo, alg geo.Algorithm) {
	for _, file := range files {

		if file.IsDir() == false && filepath.Ext(file.Name()) == ".gpx" {
			gpxDoc, err := gpxs.ParseFile(filepath.Join(gpxFilesDirectory, file.Name()), alg)
			if err != nil {
				panic(err)
			}
			var extension = filepath.Ext(file.Name())
			var name = file.Name()[0 : len(file.Name())-len(extension)]
			fileName := fmt.Sprintf("%s-graph-%s.png", name, alg.String())
			graphPath := filepath.Join(graphFilesDirectory, fileName)
			title := fmt.Sprintf("%s - %s", name, alg.String())
			createPlot(gpxDoc, title, graphPath)
		}
	}
}

func createPlot(gpx *geo.GPX, title string, filePath string) {

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(gpx.MovementStats.OverallData.Duration)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(gpx.MovementStats.MovingData.Duration)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(gpx.MovementStats.StoppedData.Duration)))
	p.Title.Text = fmt.Sprintf("%s\nDistance: %fkm - Duration: %s\nMoving Distance: %fkm - Duration: %s\nStopped Distance: %fkm - Duration: %s", title, gpx.MovementStats.OverallData.Distance/1000, t01, gpx.MovementStats.MovingData.Distance/1000, t02, gpx.MovementStats.StoppedData.Distance/1000, t03)

	p.X.Label.Text = "Distance"
	p.Y.Label.Text = "Speed (m/s)"

	points := graphPoints(gpx)

	// Make a line plotter with points and set its style.
	lpLineOverallPoints, lpPointsOverallPoints, err := plotter.NewLinePoints(points[2])
	if err != nil {
		fmt.Println(title)
		fmt.Println(err)
		return
	}

	// Make a scatter plotter and set its style.
	scatterStoppedPoints, err := plotter.NewScatter(points[3])
	if err != nil {
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
	lpLineX1, _, err := plotter.NewLinePoints(points[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	lpLineX1.Color = color.RGBA{6, 149, 171, 0xff}

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

	ptAll := make(plotter.XYs, gpx.PointsCount)
	ptStopped := make(plotter.XYs, gpx.PointsCount)

	var distance float64
	i := 0
	for _, track := range gpx.Tracks {
		for _, segment := range track.Segments {
			for x, point := range segment.Points {

				// SD
				if segment.MovementStats.SD.Valid {
					ptX1[i].X = float64(i)
					ptX1[i].Y = segment.MovementStats.SD.X1

					ptX2[i].X = float64(i)
					ptX2[i].Y = segment.MovementStats.SD.X2
				} else {
					// We assume that the SpeedThreshold is used; the value is fixed set in the customer algorithm
					ptX1[i].X = float64(i)
					ptX1[i].Y = 1.0
				}

				distance += (point.Distance / 1000)

				ptAll[i].X = float64(i)
				ptAll[i].Y = point.Speed

				if !point.IsMoving {
					ptStopped[i].X = float64(i)
					ptStopped[i].Y = point.Speed
				}
				if x == 0 {
					ptAll[0].X = 0
					ptAll[i].Y = segment.Points[x+1].Speed
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
