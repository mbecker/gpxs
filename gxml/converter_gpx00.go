package gxml

import (
	"math"
	"time"

	"github.com/mbecker/gpxs/generic"
	"github.com/mbecker/gpxs/geo"
)

/* Converter for all baseline xml elements defined in gpx00 */

//Converter00GPX00DocTracks return the geo.GPXTrack
func Converter00GPX00DocTracks(gpxDoc *geo.GPX, gpx00DocTracks []*GPX00GpxTrk, algorithm geo.Algorithm) {
	if gpx00DocTracks == nil {
		return
	}

	gpxDoc.Tracks = make([]geo.GPXTrack, len(gpx00DocTracks))
	for trackNo, track := range gpx00DocTracks {
		gpxTrack := new(geo.GPXTrack)
		// Get Track name and append it to the gpxTrackNames to return
		gpxTrack.Name = track.Name
		if len(gpxDoc.Name) == 0 {
			gpxDoc.Name = track.Name
		}

		gpxTrack.Comment = track.Cmt
		gpxTrack.Description = track.Desc
		gpxTrack.Source = track.Src
		gpxTrack.Timestamp = new(time.Time)

		// ToDo: If track.Number is not given in the xml then assign the trackNo
		if track.Number.Null() {
			gpxTrack.Number = trackNo
		} else {
			gpxTrack.Number = track.Number.Value()
		}

		gpxTrack.Type = track.Type

		if track.Segments != nil {
			gpxTrack.Segments = make([]geo.GPXTrackSegment, len(track.Segments))

			for segmentNo, segment := range track.Segments {
				gpxSegment := geo.GPXTrackSegment{}
				if segment.Points != nil {
					// Make a slice for gpxSegment.Points with the size of the slize segment.Points (xml) to store all GPXPoints
					gpxSegment.Points = make([]geo.GPXPoint, len(segment.Points))

					// Set the first point in the slice gpxSegment.Points that the second one can use this point for caluclation of distance, duration
					var prevPoint geo.GPXPoint
					gpxSegment.Points[0] = *convertPointFromGpx00(segment.Points[0])

					// Set the time of the track (if it's not already set) and segment (since it's the first point the segment does not yet have any timestamp)
					gpxSegment.Timestamp = &gpxSegment.Points[0].Timestamp
					if gpxTrack.Timestamp.IsZero() {
						gpxTrack.Timestamp = &gpxSegment.Points[0].Timestamp
					}

					for index := 1; index < len(segment.Points); index++ {
						prevPoint = gpxSegment.Points[index-1]

						// Get XML Point and convert into GPXPoint
						gpxPoint := *convertPointFromGpx00(segment.Points[index])

						// Set the point data: Duration, Distance, Speed, Pace
						gpxPoint.Point.SetPointData(prevPoint.Point, algorithm)

						// Standard deviation: Add the duration to durationSum
						gpxSegment.Duration += gpxPoint.Point.Duration
						gpxSegment.Distance += gpxPoint.Point.Distance

						// Add GPXPoint to Slice
						gpxSegment.Points[index] = gpxPoint
					}

					// Set gpxSegment data

					// Standard deviation to have a baseline of points for the calculations
					// 1. Define the mean mu (μ) for a population series: All summed values / count of values
					μ := gpxSegment.Duration / float64(len(gpxSegment.Points)) // The mean mu (μ) for a population series

					// 2.a) Define Deviation for each point: (x1−μ)
					// 2.b) Square each deviation: (x1−μ)^2
					// 2.c) Sum all squared deviation from each point
					var squaredDeviationSum float64 // Sum of all squared deviation from each point
					for index := 1; index < len(gpxSegment.Points); index++ {
						point := gpxSegment.Points[index]
						squaredDeviationSum += math.Pow(point.Duration-μ, 2)
					}

					// 3. Define the variance of the population: Divide the sum of all squared deviation of each points by the number of the population (in the previous step we used all point except the first one: len(seg.Points)-1)
					variance := squaredDeviationSum / float64((len(gpxSegment.Points) - 1))

					// 4. Define the standard deviation
					standardDeviation := math.Sqrt(variance)

					// 5. Define the the x1 and x2 value in which the points should be (sigma σ defines the range)
					x1 := μ - algorithm.Sigma()*standardDeviation
					x2 := μ + algorithm.Sigma()*standardDeviation

					// Create a slice with the length (not capacity!) of the slice gpxSegment.Points  to store all points which belongs to MovingData based on standard deviation
					gpxSegment.MovingData.Points = gpxSegment.Points[:0]                                      // https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
					gpxSegment.MovingData.Points = append(gpxSegment.MovingData.Points, gpxSegment.Points[0]) // Add the first point in the segment to the moving data since it's the starting point
					previousGPXPoint := gpxSegment.Points[0]
					gpxSegment.MovingData.MaxPace = 1000000 // ToDo: The initial vlaue is used for the if codition below (because MaxPace is basically the fastest (smallest) pace in min/km)
					for index := 1; index < len(gpxSegment.Points); index++ {
						previousGPXPoint = gpxSegment.Points[index-1]
						gpxPoint := gpxSegment.Points[index]
						if x1 <= gpxPoint.Point.Duration && gpxPoint.Point.Duration <= x2 {
							gpxPoint.Point.SetPointData(previousGPXPoint.Point, algorithm)

							gpxSegment.MovingData.MovingTime += gpxPoint.Point.Duration
							gpxSegment.MovingData.MovingDistance += gpxPoint.Point.Distance
							// Max Speed
							if gpxPoint.Point.Speed > gpxSegment.MovingData.MaxSpeed {
								gpxSegment.MovingData.MaxSpeed = gpxPoint.Point.Speed
							}
							// Max Pace
							if gpxPoint.Point.Pace > 0.0 && gpxPoint.Point.Pace < gpxSegment.MovingData.MaxPace {
								gpxSegment.MovingData.MaxPace = gpxPoint.Point.Pace
							}

							// Add GPXPoint to gpxSegment.MovingData.Points
							gpxSegment.MovingData.Points = append(gpxSegment.MovingData.Points, gpxPoint)
						} else {
							gpxSegment.MovingData.StoppedTime += gpxPoint.Point.Duration
							gpxSegment.MovingData.StoppedDistance += gpxPoint.Point.Distance
						}

					}
					gpxSegment.MovingData.Duration += gpxSegment.MovingData.MovingTime + gpxSegment.MovingData.StoppedTime
					gpxSegment.MovingData.Distance += gpxSegment.MovingData.MovingDistance + gpxSegment.MovingData.StoppedDistance

					averageSpeed, errAverageSpeed := algorithm.Speed(gpxSegment.MovingData.MovingDistance, gpxSegment.MovingData.MovingTime)
					if errAverageSpeed == nil {
						gpxSegment.MovingData.AverageSpeed = averageSpeed
					}
					averagePace, errAveragePace := algorithm.Pace(gpxSegment.MovingData.MovingDistance, gpxSegment.MovingData.MovingTime)
					if errAveragePace == nil {
						gpxSegment.MovingData.AveragePace = averagePace
					}

					// Append gpxSegent to gpxTrackSegments
					gpxTrack.Segments[segmentNo] = gpxSegment

					// gpxTrack MovinngData
					gpxTrack.MovingData.SetMovingValues(
						gpxSegment.MovingData.MovingTime,
						gpxSegment.MovingData.StoppedTime,
						gpxSegment.MovingData.MovingDistance,
						gpxSegment.MovingData.StoppedDistance,
						gpxSegment.MovingData.MaxSpeed,
						gpxSegment.MovingData.AverageSpeed,
						gpxSegment.MovingData.MaxPace,
						gpxSegment.MovingData.AveragePace,
						segmentNo+1,
					)

				}
			}

		}

		gpxDoc.Tracks[trackNo] = *gpxTrack

		// Set the time of the gpxDoc (if it's not already set)
		if gpxDoc.Timestamp.IsZero() {
			gpxDoc.Timestamp = gpxTrack.Timestamp
		}

		// gpxDoc MovinngData
		gpxDoc.MovingData.SetMovingValues(
			gpxTrack.MovingData.MovingTime,
			gpxTrack.MovingData.StoppedTime,
			gpxTrack.MovingData.MovingDistance,
			gpxTrack.MovingData.StoppedDistance,
			gpxTrack.MovingData.MaxSpeed,
			gpxTrack.MovingData.AverageSpeed,
			gpxTrack.MovingData.MaxPace,
			gpxTrack.MovingData.AveragePace,
			trackNo+1,
		)
	}

}

// Set00GPX00DocWaypoint sets the gpxDoc.Waypoint if the xml has points (GPX00GpxPoint)
func Set00GPX00DocWaypoint(gpxDoc *geo.GPX, gpx00Waypoints []*GPX00GpxPoint) {
	if gpx00Waypoints != nil {
		waypoints := make([]geo.GPXPoint, len(gpx00Waypoints))
		for waypointNo, waypoint := range gpx00Waypoints {
			waypoints[waypointNo] = *convertPointFromGpx00(waypoint)
		}
		gpxDoc.Waypoints = waypoints
	}
}

// Set00GPX00DocRoutes sets the gpxDoc.Routes if the xml has routes (GPX00GpxRte)
func Set00GPX00DocRoutes(gpxDoc *geo.GPX, gpx00Rte []*GPX00GpxRte) {
	if gpx00Rte != nil {
		gpxDoc.Routes = make([]geo.GPXRoute, len(gpx00Rte))
		for routeNo, route := range gpx00Rte {
			r := new(geo.GPXRoute)

			r.Name = route.Name
			r.Comment = route.Cmt
			r.Description = route.Desc
			r.Source = route.Src
			// TODO
			//r.Links = route.Links
			if route.Number.Null() {
				r.Number = routeNo
			} else {
				r.Number = route.Number.Value()
			}
			r.Type = route.Type
			// TODO
			//r.RoutePoints = route.RoutePoints

			if route.Points != nil {
				r.Points = make([]geo.GPXPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = *convertPointFromGpx00(point)
				}
			}

			gpxDoc.Routes[routeNo] = *r
		}
	}
}

func convertPointFromGpx00(original *GPX00GpxPoint) *geo.GPXPoint {
	result := new(geo.GPXPoint)
	result.Latitude = original.Lat
	result.Longitude = original.Lon
	result.Elevation = original.Ele
	time, _ := parseGPXTime(original.Timestamp)
	if time != nil {
		result.Timestamp = *time
	}
	result.MagneticVariation = original.MagVar
	result.GeoidHeight = original.GeoIdHeight
	result.Name = original.Name
	result.Comment = original.Cmt
	result.Description = original.Desc
	result.Source = original.Src
	// TODO
	//w.Links = original.Links
	result.Symbol = original.Sym
	result.Type = original.Type
	result.TypeOfGpsFix = original.Fix
	if original.Sat != nil {
		result.Satellites = *generic.NewNullableInt(*original.Sat)
	}
	if original.Hdop != nil {
		result.HorizontalDilution = *generic.NewNullableFloat64(*original.Hdop)
	}
	if original.Vdop != nil {
		result.VerticalDilution = *generic.NewNullableFloat64(*original.Vdop)
	}
	if original.Pdop != nil {
		result.PositionalDilution = *generic.NewNullableFloat64(*original.Pdop)
	}
	if original.AgeOfDGpsData != nil {
		result.AgeOfDGpsData = *generic.NewNullableFloat64(*original.AgeOfDGpsData)
	}
	if original.DGpsId != nil {
		result.DGpsId = *generic.NewNullableInt(*original.DGpsId)
	}
	return result
}

func convertPointToGpx00(original *geo.GPXPoint) *GPX00GpxPoint {
	result := new(GPX00GpxPoint)
	result.Lat = original.Latitude
	result.Lon = original.Longitude
	result.Ele = original.Elevation
	result.Timestamp = formatGPXTime(&original.Timestamp)
	result.MagVar = original.MagneticVariation
	result.GeoIdHeight = original.GeoidHeight
	result.Name = original.Name
	result.Cmt = original.Comment
	result.Desc = original.Description
	result.Src = original.Source
	// TODO
	//w.Links = original.Links
	result.Sym = original.Symbol
	result.Type = original.Type
	result.Fix = original.TypeOfGpsFix
	if original.Satellites.NotNull() {
		value := original.Satellites.Value()
		result.Sat = &value
	}
	if original.HorizontalDilution.NotNull() {
		value := original.HorizontalDilution.Value()
		result.Hdop = &value
	}
	if original.VerticalDilution.NotNull() {
		value := original.VerticalDilution.Value()
		result.Vdop = &value
	}
	if original.PositionalDilution.NotNull() {
		value := original.PositionalDilution.Value()
		result.Pdop = &value
	}
	if original.AgeOfDGpsData.NotNull() {
		value := original.AgeOfDGpsData.Value()
		result.AgeOfDGpsData = &value
	}
	if original.DGpsId.NotNull() {
		value := original.DGpsId.Value()
		result.DGpsId = &value
	}
	return result
}
