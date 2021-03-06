# gpxs - gpxgo on steriods

> All credits to gpx gpxgo - https://github.com/tkrajina/gpxgo/issues

# Motivation
I'm using Strava quite heavily (3-5 activities/week) for most of my sport activities like running, cycling, hiking, etc. (see https://www.strava.com/athletes/8844168). At the moment I'm using the Strava iOS App.
Unfortunatley I started with Runkeeper years ago and switched for a time to Garmin and other tools. Once I exported (more or less all) my tracks from all third parties and uploaded the gpx files to Strava. Strava analyzed all the files for me and added the activities to my statistics. Nice!

I want to have all my activites on my local desktop to analyze the gpx files and create nice stats and graphics. For that I downloaded all my gpx files from Strava and started to look how to parse the gpx files. That's how I discovered gpxgo.

(How to download all activities see: https://support.strava.com/hc/en-us/articles/216918437-Exporting-your-Data-and-Bulk-Export #Section: Bulk Export at the the end of the page)

The fantastic library gpxgo (https://github.com/tkrajina/gpxgo/issues) by tkrajina parses gpx files and returns information of the gpx/tracks/segements/points like distance, duration and movingtime/distance + stoppedtime/distance.

Unfortunatley the gpx information like distance (movingdistance), duration (movingtime) does not match with the Strava information.

My assumption is that Strava does a quite good in analyzing gpx data and should be the baseline to match. I looked into the libraray and noticed for example that the distance calculation is quite fast but less accurate than other GPS distance calculation formulas like Haversine and Vincenty. I started to implement these formulas in the libraray and quickly noticed that for each part of the gpx (gpx->track->segment->point) the data is re-calculated new. Additionaly the points which are used for the caluclation for "Moving"Time/Distance are based on the assumption that the speed between two points must be greater than 1.0m/s. That results in a big difference of "Moving"Time/Distance of Strava and the library.

**So I've started to update te libraray as follows:**

- Change the method how the data is aggregated (from point->segments->tracks->gpx)
- Refactoring of the code into multiple packages; refactoring the XML strucst/converter since the gpx v1.0 and v1.1 isn't so different
- Use track name as gpx name if it's not already set; Use first point timestamp for track/gpx timestamp if it's not already set
- User provides the the caluclation of distance, duration, speed, pace, etc. + method to normalize the points
- Added new GPS distance cluclation
- Use standard deviation to determine which points should be used for "Moving"Time/Distance

**How is the data now agregated?**

At first the point data like duration, distane between the previous and the actual point is caluclated. The the data is then aggregated to the next level as follows:

```
1.) The Point data is parsed from the XML and additional parameters are calculated as defined the 'type Point struct':
type Point struct {
	Latitude  float64
	Longitude float64
	Elevation generic.NullableFloat64
	Timestamp time.Time

	Distance float64 // The distance (m) from the previos to this point
	Duration float64 // The duration (sec) from the previos to this point
	Speed    float64 // The speed (m/s) from the previous point to this point
	Pace     float64 // The pace (m/s) from the previous point to this point
}

2.) The sum of all Points is used to define the Segment's data as defined in the 'type GPXTrackSegment struct' field 'MovingData':

type GPXTrackSegment struct {
	Points     []GPXPoint
	MovingData MovingData
}

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

3.) The sum of all Segements is used to define the Track's data (see 2.)
//GPXTrack implements a gpx track
type GPXTrack struct {
	Segments   []GPXTrackSegment
	MovingData MovingData
}

4.) The sum of all Track is used to define GPX's data (see 2.)
type GPX struct {
	MovingData MovingData
	Tracks    []GPXTrack
}
```

**Refactoring**

For my personal style to work with multiple structs, interface and funcs I splitted the code into multiple packages. That makes it easier for me to work with and to understand the code.
The structs which define the XML structure of gpx v1.0 and v1.0 are quite similar. So I tried to use structs as a baseline (call it v00) and then only implement structs for the specifc XML elements which differs. With that the conversion from XML to GPX (and vice versa) is only implemented once (and not twice).

**GPX name and timestamp from track and first point**

In the gpx files from Strava (and the other tird parties I've used) no name and sometimes no timestamp is given for gpx at the top level. If that's the case then __gpxs__ uses the name of the track(s) and the timestamp of the first point. The name and timestamp is taken by FIFO (first in first out); that means basically that first name/timestamp is used.

**User method to calculate distance, duration, speed, pace + method to normalize the points**

The idea is that the users who use this libraray provides own methods to calculate the points data. That means the users could use a GPS distance calculatio not yet provided by the library. Additionaly you could use your own standardize method to use only the points you want to use for your calculation.

Methods implemented:
- "Standard Method" by gpxgo: Basic but fast calculation of distance between two points
- Vincenty: More accurate methods
- Normalize method (1): Threshold by gpxgo; use only points above speed threshold of 1m/s
- Normalize method (2): Standard deviation (see below)

**New GPS distance cluclation**

Added a method to use GPS formula Vincenty. Also added the calculation methods from gpxgo as the "Standard Algorithm".

**Standard deviation for "Moving"Time/Distance?**

A point has the following data: lat/long (where) and timestamp (when). With the information where and when of each point you could calculate the distance (how long in meters?) and duration (how long in seconds?) between the previous and actual point.

It's quite normal that you stop during an activity (waiting on the street, doing a longer pause, etc) and then the distance/duration between the point before and after the pause could be quite long. Let's say your are doing a cycling tour to a castle; at the castle you stop (Point p0(position0, time0)) and after 1hour you are cycling back (start of the return Point p1(position0+100m, timee0+3600seconds)). The speed would be 

```
P1(position0+100m)-P0(position0) / P1(time0+36000) - P0(time0) -> 100m / 3600sec -> ~0,0278 m/s = 0,10008 km/h
```

Obvious the data between these two point p0 and p1 shouldn't be used for the "Moving"Time/Distance since p1 is a new startig point. The point after P1 should then be used for the "Moving"Time/Distance.

The idea of gpxgo is that the speed between two points must be grater than 1m/s; then the points are used for the caluclation. The approach is fair enough since it's quite fast. But for me the different to Strava is too much (for example the duration for a cycling tour with gpxgo is 3h20min and with Stava it's 2h45min which is the "real" movingtime from my trip).

The standard deviation (https://en.wikipedia.org/wiki/Standard_deviation) uses all points data of the gpx track and "quantifies the amount of variation or dispersion of the set of data values". In my words: You define a weightehd average of all points duration and then define that only the points should be used which are close enough to that duration average (like 90% or 95%).
With an standard deviation Sigma σ of 1.644854 (~95%) I'm quite close to the Strava information (only some seconds in each of my ~1050 gpx files).

So good enough for me!

**Performance**

See benchmarks below.

My assumption was that the "new" libraray gpxs is faster than gpxgo because of the new aggregation method. For a benchmark of Vincenty I added the method to gpxgo. gpxgo does not use the standard deviation; so __gpxgs__ is instrumented to use the same "normalization method" as gpxgo.


__**Conclusion**__

The reuslts of gpxs (distance, "moving"time/distance, "stopped"time/distance) are for **me more accurate**. "More accurate" **means for me** the calculated infomation is closer to Stava which is **my baseline**.

The "more accurate" is achieved by using a new GPS distance calculation (Vincenty formula) and the standard deviation to nromalize the data set (which points should be used for the "moving"time/distance calculation).

Suprisingly the new aggregation is not so faster as I thought ("only" 2seconds); using the standard deviation it's faster (4sec).

Additionaly the new approach to use custom methods to normalize the data sete and caluclate the information provides more flexibility.


# TODO:

- [ ]  Check that csv files are generated correctly
- [ ]  Check what gpx.Routes and gpx.Waypoint are and how to use it in gpxs
- [ ]  Test that GPX Duration - Sum of all Track Duration - Sum of all Tracks Segments Duration - Sm of all Tracks Segments MovingData Duration are equal
- [ ] Write test files
- [O] Test: Sum of all Track Distance - Sum of all Tracks Segments Distance - Sm of all Tracks Segments MovingData Distance are equal
- [ ] Write benchmark fies

# How to use / Examples

1.) Built-in geo.Algorithm
```go
// examples/example1.go

package examples

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
)

func example1() {
	// 1.) Use a built-in geo.Algorithm
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

	// 2.) Parse a gpx file with the geo.Algorithm
	gpxDoc, err := gpxs.ParseFile(filepath.Join("fileDirectory", "test.gpx"), &vincenty)
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

	t01, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(duration)))
	t02, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(movingTime)))
	t03, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(stoppedTime)))

	t04, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackDuration)))
	t05, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackMovingTime)))
	t06, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(trackStoppedTime)))

	t07, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentDuration)))
	t08, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentMovingTime)))
	t09, _ := time.ParseDuration(fmt.Sprintf("%ds", int64(segmentStoppedTime)))

	fmt.Printf("------")
	fmt.Printf("%s", t01)                    // GPX Duration
	fmt.Printf("%f", distance/1000.0)        // GPX Distance
	fmt.Printf("%s", t02)                    // GPX Moving Time)
	fmt.Printf("%s", t03)                    // GPX Stopped Time
	fmt.Printf("%f", movingDistance/1000.0)  // GPX Moving Distance
	fmt.Printf("%f", stoppedDistance/1000.0) // GPX Stopped Distance
	fmt.Printf("------")
	fmt.Printf("%s", t04)                         // Track Duration
	fmt.Printf("%f", trackDistance/1000.0)        // Track Moving Distance
	fmt.Printf("%s", t05)                         // Track Moving Time
	fmt.Printf("%s", t06)                         // Track Stopped Time
	fmt.Printf("%f", trackMovingDistance/1000.0)  // Track Moving Distance
	fmt.Printf("%f", trackStoppedDistance/1000.0) // Track Stopped Distance
	fmt.Printf("------")
	fmt.Printf("%s", t07)                           // Segment Duration
	fmt.Printf("%f", segmentDistance/1000.0)        // Segment Distance
	fmt.Printf("%s", t08)                           // Segment Moving Time
	fmt.Printf("%s", t09)                           // Segment  Stopped Time
	fmt.Printf("%f", segmentMovingDistance/1000.0)  // Segment Moving Distance
	fmt.Printf("%f", segmentStoppedDistance/1000.0) // Segment  Stopped Distance
}
```
2.) Custom geo.Algorithm
```go
// examples/example2.go

package exmaples

import (
	"errors"
	"path/filepath"

	"github.com/mbecker/gpxs/geo"
	"github.com/mbecker/gpxs/gpxs"
)

// 1.) Define customer geo.algorithm; must fullfill interface geo.Algorithm
type CustomAlgorithm struct {
	CustomParameter float64
}

// ShouldStandardDeviation (CustomAlgorithm) returns if the standard deviation should be used or not
func (c *CustomAlgorithm) ShouldStandardDeviation() bool {
	return false
}

// Sigma (CustomAlgorithm) returns the sigma for the standard deviation; not used because of 'ShouldStandardDeviation' returns 'false'
func (c *CustomAlgorithm) Sigma() float64 {
	return 0
}

// Duration (CustomAlgorithm) returns the time.Duration from point p1 to previousPoint in sec
func (c *CustomAlgorithm) Duration(p1 *geo.Point, previousPoint *geo.Point) (float64, error) {
	return p1.Timestamp.Sub(previousPoint.Timestamp).Seconds(), nil
}

// CustomMovingPoints (CustomAlgorithm) defines which points should be used for "Moving"Time/Distance and if the it's set the new gpxPoint.Point Data
func (c *CustomAlgorithm) CustomMovingPoints(gpxPoint *geo.GPXPoint, previousGPXPoint *geo.GPXPoint, algorithm geo.Algorithm) error {

	/* 	Define which points should be used; if a point should be used for calculation then set it's new values like Duration, Distance, Speed, etc.
	Here we use the set the new value for the points which used for "Moving"Time/Distanc
	*/

	// speed < 100 m/s
	if gpxPoint.Speed < 100.0 {
		return errors.New("Point Speed below threshold")
	}
	gpxPoint.Point.SetPointData(previousGPXPoint.Point, algorithm)
	return nil
}

// Distance (CustomAlgorithm) returns just 100 as an example
func (c *CustomAlgorithm) Distance(p1 *geo.Point, previousPoint *geo.Point) (float64, error) {
	return 101, nil
}

// Speed (CustomAlgorithm) returns the speed in m/s
func (c *CustomAlgorithm) Speed(distance float64, duration float64) (float64, error) {
	return 101.9, nil
}

// Pace (CustomAlgorithm) returns the pace in s/m
func (c *CustomAlgorithm) Pace(distance float64, duration float64) (float64, error) {
	return 20.9, nil
}

func example2() {

	customAlgorithm := CustomAlgorithm{
		CustomParameter: 100.9,
	}

	// 2.) Parse a gpx file with the custom geo.Algorithm
	gpxDoc, err := gpxs.ParseFile(filepath.Join("fileDirectory", "test.gpx"), &customAlgorithm)
	if err != nil {
		panic(err)
	}

	// 3.) See how to use gpxDoc in the first example

}
```



# Bechmark: Different gps distance calculations

> See Google Sheets: https://docs.google.com/spreadsheets/d/1NapqlDL7DBYxzw53JUSJokM6-iCYHhJsiehCRhK2kDw/edit?usp=sharing

**gpxs**
> Benchmark with different GPS distance calculation methods; different normalization methods (standard deviation vs. default speed threshold)

|           TYPE           |     VINCENTY W/O SD            |     VINCENTY WITH SD           |
|--------------------------|--------------------------------|--------------------------------|
| # of files               |                           1047 |                           1047 |
| ------                   | ------                         | ------                         |
| GPX Duration             | 1009h13m1s                     | 1009h13m1s                     |
| GPX Distance             |                   11543.863484 |                   11543.863484 |
| GPX Moving Time          | 931h54m37s                     | 894h42m47s                     |
| GPX Stopped Time         | 77h18m24s                      | 114h30m14s                     |
| GPX Moving Distance      |                   11462.968677 |                   10886.758060 |
| GPX Stopped Distance     |                      80.894806 |                     657.105423 |
| ------                   | ------                         | ------                         |
| Track Duration           | 1009h13m1s                     | 1009h13m1s                     |
| Track Distance           |                   11543.863484 |                   11543.863484 |
| Track Moving Time        | 931h54m37s                     | 894h42m47s                     |
| Track Stopped Time       | 77h18m24s                      | 114h30m14s                     |
| Track Moving Distance    |                   11462.968677 |                   10886.758060 |
| Track Stopped Distance   |                      80.894806 |                     657.105423 |
| ------                   | ------                         | ------                         |
| Segment Duration         | 1009h13m1s                     | 1009h13m1s                     |
| Segment Distance         |                   11543.863484 |                   11543.863484 |
| Segment Moving Time      | 931h54m37s                     | 894h42m47s                     |
| Segment Stopped Time     | 77h18m24s                      | 114h30m14s                     |
| Segment Moving Distance  |                   11462.968677 |                   10886.758060 |
| Segment Stopped Distance |                      80.894806 |                     657.105423 |
| ------                   | ------                         | ------                         |
| Execution time           | 21.2726596s                    | 21.7201131s                    |


|           TYPE           |    STANDARD (LENGTH2D) W/O SD  |    STANDARD (LENGTH2D) WITH SD |
|--------------------------|--------------------------------|--------------------------------|
| # of files               |                           1047 |                           1047 |
| ------                   | ------                         | ------                         |
| GPX Duration             | 1009h13m1s                     | 1009h13m1s                     |
| GPX Distance             |                   11515.996154 |                   11515.996154 |
| GPX Moving Time          | 931h50m40s                     | 894h42m47s                     |
| GPX Stopped Time         | 77h22m21s                      | 114h30m14s                     |
| GPX Moving Distance      |                   11435.016290 |                   10860.516895 |
| GPX Stopped Distance     |                      80.979864 |                     655.479259 |
| ------                   | ------                         | ------                         |
| Track Duration           | 1009h13m1s                     | 1009h13m1s                     |
| Track Distance           |                   11515.996154 |                   11515.996154 |
| Track Moving Time        | 931h50m40s                     | 894h42m47s                     |
| Track Stopped Time       | 77h22m21s                      | 114h30m14s                     |
| Track Moving Distance    |                   11435.016290 |                   10860.516895 |
| Track Stopped Distance   |                      80.979864 |                     655.479259 |
| ------                   | ------                         | ------                         |
| Segment Duration         | 1009h13m1s                     | 1009h13m1s                     |
| Segment Distance         |                   11515.996154 |                   11515.996154 |
| Segment Moving Time      | 931h50m40s                     | 894h42m47s                     |
| Segment Stopped Time     | 77h22m21s                      | 114h30m14s                     |
| Segment Moving Distance  |                   11435.016290 |                   10860.516895 |
| Segment Stopped Distance |                      80.979864 |                     655.479259 |
| ------                   | ------                         | ------                         |
| Execution time           | 16.6517713s                    | 17.1150661s                    |


|           TYPE           |    STANDARD (LENGTH3D) W/O SD  |    STANDARD (LENGTH3D) WITH SD |
|--------------------------|--------------------------------|--------------------------------|
| # of files               |                           1047 |                           1047 |
| ------                   | ------                         | ------                         |
| GPX Duration             | 1009h13m1s                     | 1009h13m1s                     |
| GPX Distance             |                   11539.633336 |                   11539.633336 |
| GPX Moving Time          | 932h20m46s                     | 894h42m47s                     |
| GPX Stopped Time         | 76h52m15s                      | 114h30m14s                     |
| GPX Moving Distance      |                   11457.471258 |                   10882.709172 |
| GPX Stopped Distance     |                      82.162079 |                     656.924164 |
| ------                   | ------                         | ------                         |
| Track Duration           | 1009h13m1s                     | 1009h13m1s                     |
| Track Distance           |                   11539.633336 |                   11539.633336 |
| Track Moving Time        | 932h20m46s                     | 894h42m47s                     |
| Track Stopped Time       | 76h52m15s                      | 114h30m14s                     |
| Track Moving Distance    |                   11457.471258 |                   10882.709172 |
| Track Stopped Distance   |                      82.162079 |                     656.924164 |
| ------                   | ------                         | ------                         |
| Segment Duration         | 1009h13m1s                     | 1009h13m1s                     |
| Segment Distance         |                   11539.633336 |                   11539.633336 |
| Segment Moving Time      | 932h20m46s                     | 894h42m47s                     |
| Segment Stopped Time     | 76h52m15s                      | 114h30m14s                     |
| Segment Moving Distance  |                   11457.471258 |                   10882.709172 |
| Segment Stopped Distance |                      82.162079 |                     656.924164 |
| ------                   | ------                         | ------                         |
| Execution time           | 17.3354035s                    | 17.0341751s                    |


> gpxgo Benchmark - The execution with the Vincenty forula takes longer due to more caluclation steps

## gpxgo 
> Simple distance calculation (length2D, length3D), normalization of poins with default speed threshold (1m/s)
> Vincenty formula added by me

|           TYPE           |   LENGTH2D   |   LENGTH3D   |   VINCENTY   |
|--------------------------|--------------|--------------|--------------|
| # of files               |         1047 |         1047 |         1047 |
| ------                   | ------       | ------       | ------       |
| GPX Duration             | 1009h13m1s   | 1009h13m1s   | 1009h13m1s   |
| GPX Distance             | 11515.996154 | 11539.633336 |    11.543863 |
| GPX Moving Time          | 963h28m41s   | 963h28m41s   | 963h28m41s   |
| GPX Stopped Time         | 45h44m20s    | 45h44m20s    | 45h44m20s    |
| GPX Moving Distance      | 11523.183462 | 11523.183462 | 11523.183462 |
| GPX Stopped Distance     |    16.449875 |    16.449875 |    16.449875 |
| ------                   | ------       | ------       | ------       |
| Track Duration           | 1009h13m1s   | 1009h13m1s   | 1009h13m1s   |
| Track Distance           | 11515.996154 | 11539.633336 |    11.543863 |
| Track Moving Time        | 963h28m41s   | 963h28m41s   | 963h28m41s   |
| Track Stopped Time       | 45h44m20s    | 45h44m20s    | 45h44m20s    |
| Track Moving Distance    | 11523.183462 | 11523.183462 | 11523.183462 |
| Track Stopped Distance   |    16.449875 |    16.449875 |    16.449875 |
| ------                   | ------       | ------       | ------       |
| Segment Duration         | 1009h13m1s   | 1009h13m1s   | 1009h13m1s   |
| Segment Distance         |     0.000000 |     0.000000 |     0.000000 |
| Segment Moving Time      | 963h28m41s   | 963h28m41s   | 963h28m41s   |
| Segment Stopped Time     | 45h44m20s    | 45h44m20s    | 45h44m20s    |
| Segment Moving Distance  | 11523.183462 | 11523.183462 | 11523.183462 |
| Segment Stopped Distance |    16.449875 |    16.449875 |    16.449875 |
| ------                   | ------       | ------       | ------       |
| Execution time           | 17.9543938s  | 18.4098755s  | 25.3006466s  |

> gpxgo Benchmark - The execution with the Vincenty forula takes longer due to more caluclation steps

# Standard Deviation

See https://en.wikipedia.org/wiki/Standard_deviation
```
3.29053 	~ 99.9%
2.57583 	~ 99%
2.17009 	~ 97%
1.959964 	~ 95%
```

> From analysis see graphs in '/test/graphs' the Sigma 3.29053 (~99.9%) seems to fit best

|  **Confidence <br/>interval** | **Proportion within** | **Proportion without** |  |
|  :------: | :------: | :------: | ------ |
|   | **Percentage** | **Percentage** | **Fraction** |
|  0.318 639σ | 0.25 | 0.75 | 3 / 4 |
|  0.674490σ | 0.5 | 0.5 | 1 / 2 |
|  0.994458σ | 0.68 | 0.32 | 1 / 3.125 |
|  1σ | 0.682689492 | 0.317310508 | 1 / 3.1514872 |
|  1.281552σ | 0.8 | 0.2 | 1 / 5 |
|  1.644854σ | 0.9 | 0.1 | 1 / 10 |
|  1.959964σ | 0.95 | 0.05 | 1 / 20 |
|  2σ | 0.954499736 | 0.045500264 | 1 / 21.977895 |
|  2.575829σ | 0.99 | 0.01 | 1 / 100 |
|  3σ | 0.997300204 | 0.002699796 | 1 / 370.398 |
|  3.290527σ | 0.999 | 0.001 | 1 / 1000 |
|  3.890592σ | 0.9999 | 0.0001 | 1 / 10000 |
|  4σ | 0.99993666 | 0.00006334 | 1 / 15787 |
|  4.417173σ | 0.99999 | 0.00001 | 1 / 100000 |
|  4.5σ | 0.999993204653751 | 0.000006795346249 | 3.4 / 1000000 (on each side of mean) |
|  4.891638σ | 0.999999 | 0.000001 | 1 / 1000000 |
|  5σ | 0.999999426697 | 5.73303e-7 | 1 / 1744278 |
|  5.326724σ | 0.9999999 | 1e-7 | 1 / 10000000 |
|  5.730729σ | 0.99999999 | 1e-8 | 1 / 100000000 |
|  6σ | 0.999999998027 | 1.973e-9 | 1 / 506797346 |
|  6.109410σ | 0.999999999 | 1e-9 | 1 / 1000000000 |
|  6.466951σ | 0.9999999999 | 1e-10 | 1 / 10000000000 |
|  6.806502σ | 0.99999999999 | 1e-11 | 1 / 100000000000 |
|  7σ | 0.99999999999744 | 2.56e-12 | 1 / 390682215445 |
