package gxml

import (
	"math"
	"strings"

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
	gpxDoc.MovementStats = geo.MovementStats{
		OverallData: geo.MovementData{},
		MovingData:  geo.MovementData{},
		StoppedData: geo.MovementData{},
	}
	for trackNo, track := range gpx00DocTracks {
		gpxTrack := new(geo.GPXTrack)
		gpxTrack.MovementStats = geo.MovementStats{
			OverallData: geo.MovementData{},
			MovingData:  geo.MovementData{},
			StoppedData: geo.MovementData{},
		}

		// Get Track name and append it to the gpxTrackNames to return
		gpxTrack.Name = track.Name
		if len(gpxDoc.Name) == 0 {
			gpxDoc.Name = track.Name
		}

		gpxTrack.Comment = track.Cmt
		gpxTrack.Description = track.Desc
		gpxTrack.Source = track.Src

		// TODO: If track.Number is not given in the xml then assign the trackNo
		if track.Number.Null() {
			gpxTrack.Number = trackNo
		} else {
			gpxTrack.Number = track.Number.Value()
		}

		// Check Type of Activity: Strava defines the activity type with a nuber ("1" == Cycling, "4" = Hiking, "9" == Running,); other parties like Garmin / Runkeeper has the activity tpye as a descriptive text in the track.Name
		gpxTrack.Type = track.Type
		if len(gpxTrack.Type) == 0 {
			activityName := strings.ToLower(track.Name)
			activityType, err := algorithm.CheckActivityType(activityName)
			if err == nil {
				gpxTrack.Type = activityType
			}
		}
		if len(gpxDoc.Type) == 0 {
			gpxDoc.Type = gpxTrack.Type
		}

		if track.Segments != nil {
			gpxTrack.Segments = make([]geo.GPXTrackSegment, len(track.Segments))

			for segmentNo, segment := range track.Segments {
				gpxSegment := geo.GPXTrackSegment{}
				gpxSegment.MovementStats = geo.MovementStats{
					OverallData: geo.MovementData{},
					MovingData:  geo.MovementData{},
					StoppedData: geo.MovementData{},
				}

				if segment.Points != nil {
					// Make a slice for gpxSegment.Points with the size of the slize segment.Points (xml) to store all GPXPoints
					gpxSegment.Points = make([]geo.GPXPoint, len(segment.Points))

					// Set the first point in the slice gpxSegment.Points that the second one can use this point for caluclation of distance, duration
					var prevPoint geo.GPXPoint
					gpxSegment.Points[0] = *convertPointFromGpx00(segment.Points[0])

					// Set the time for gpxDoc, track, segment
					gpxDoc.MovementStats.OverallData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)
					gpxDoc.MovementStats.MovingData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)

					gpxTrack.MovementStats.OverallData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)
					gpxTrack.MovementStats.MovingData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)

					gpxSegment.MovementStats.OverallData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)
					gpxSegment.MovementStats.MovingData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)

					// Define overallDuration for Standard Deviation (maybe useful for other algorthm as well?)
					var overallDuration float64

					// Loop all points and set the data of each point like duration, distance, speed, etc.
					for index := 1; index < len(segment.Points); index++ {
						prevPoint = gpxSegment.Points[index-1]

						// Get XML Point and convert into GPXPoint
						gpxPoint := *convertPointFromGpx00(segment.Points[index])

						// Set the point data: Duration, Distance, Speed, Pace
						gpxPoint.Point.SetPointData(&prevPoint.Point, algorithm)

						// ! Important for Standard deviation ! Add the duration / distance to the overall duration to know how long the duration / distance for the all points are (All point == Segment OverallData)
						overallDuration += gpxPoint.Point.Distance / gpxPoint.Point.Duration

						// Add GPXPoint to Slice
						gpxSegment.Points[index] = gpxPoint

					}

					// Create a slice with the length (not capacity!) of the slice gpxSegment.Points to save the pointer to the memory where the the GPXPoint is stored
					gpxSegment.MovementStats.MovingData.Points = make([]*geo.GPXPoint, len(segment.Points))
					gpxSegment.MovementStats.MovingData.Points = append(gpxSegment.MovementStats.MovingData.Points, &gpxSegment.Points[0])
					gpxSegment.MovementStats.MovingData.StartTime.SetTime(gpxSegment.Points[0].Timestamp.Time)

					gpxSegment.MovementStats.StoppedData.Points = make([]*geo.GPXPoint, len(segment.Points))

					var previousGPXPoint geo.GPXPoint

					if algorithm.ShouldStandardDeviation() {
						/**
							STANDARD DEVIATION -- BEGIN
						**/

						// Standard deviation to have a baseline of points for the calculations
						// 1. Define the mean mu (μ) for a population series: All summed values / count of values
						μ := overallDuration / float64(len(gpxSegment.Points)) // The mean mu (μ) for a population series

						// 2.a) Define Deviation for each point: (x1−μ)
						// 2.b) Square each deviation: (x1−μ)^2
						// 2.c) Sum all squared deviation from each point
						var squaredDeviationSum float64 // Sum of all squared deviation from each point
						for index := 1; index < len(gpxSegment.Points); index++ {
							point := gpxSegment.Points[index]
							squaredDeviationSum += math.Pow((point.Distance/point.Duration)-μ, 2)
						}

						// 3. Define the variance of the population: Divide the sum of all squared deviation of each points by the number of the population (in the previous step we used all point except the first one: len(seg.Points)-1)
						variance := squaredDeviationSum / float64((len(gpxSegment.Points) - 1))

						// 4. Define the standard deviation
						standardDeviation := math.Sqrt(variance)

						// 5. Define the the x1 and x2 value in which the points should be (sigma σ defines the range)
						x1 := μ - algorithm.Sigma()*standardDeviation
						x2 := μ + algorithm.Sigma()*standardDeviation

						// Filter all points which belongs to the standard deviation area
						for index := 1; index < len(gpxSegment.Points); index++ {
							previousGPXPoint = gpxSegment.Points[index-1]
							gpxPoint := gpxSegment.Points[index]
							gpxPoint.Point.SetPointData(&previousGPXPoint.Point, algorithm)

							pointXValue := gpxPoint.Distance / gpxPoint.Duration
							if x1 <= pointXValue && pointXValue <= x2 {
								// Points are in standard deviation area
								gpxSegment.MovementStats.MovingData.SetValues(&gpxPoint, &prevPoint, index, algorithm)
							} else {
								// Points are not in standard deviation area
								gpxSegment.MovementStats.StoppedData.SetValues(&gpxPoint, &prevPoint, index, algorithm)
							}
							// fmt.Println("2. Add to overall data")
							gpxSegment.MovementStats.OverallData.SetValues(&gpxPoint, &prevPoint, index, algorithm)

						}

					} else {
						/**
							Do not use Standard Deviation - Begin
						**/
						for index := 1; index < len(gpxSegment.Points); index++ {
							previousGPXPoint = gpxSegment.Points[index-1]
							gpxPoint := gpxSegment.Points[index]

							// Custom Moving Points
							err := algorithm.CustomMovingPoints(&gpxPoint, &previousGPXPoint, algorithm)
							if err != nil {
								// Error says: Do not use the point for "Moving"Time and "MovingDistance"
								gpxSegment.MovementStats.StoppedData.SetValues(&gpxPoint, &prevPoint, index, algorithm)
							} else {
								// TODO: the gpxPoint Data should be set by algorithm.CustomMovingPoints

								gpxSegment.MovementStats.MovingData.SetValues(&gpxPoint, &prevPoint, index, algorithm)
							}
							gpxSegment.MovementStats.OverallData.SetValues(&gpxPoint, &prevPoint, index, algorithm)

						}
					}

					// Append gpxSegent to gpxTrackSegments
					gpxTrack.Segments[segmentNo] = gpxSegment

					gpxTrack.MovementStats.OverallData.SetValuesFromMovementData(&gpxSegment.MovementStats.OverallData, segmentNo, algorithm)
					gpxTrack.MovementStats.MovingData.SetValuesFromMovementData(&gpxSegment.MovementStats.MovingData, segmentNo, algorithm)
					gpxTrack.MovementStats.StoppedData.SetValuesFromMovementData(&gpxSegment.MovementStats.StoppedData, segmentNo, algorithm)

				}
			}

		}

		gpxDoc.Tracks[trackNo] = *gpxTrack

		// gpxDoc MovinngData
		gpxDoc.MovementStats.OverallData.SetValuesFromMovementData(&gpxTrack.MovementStats.OverallData, trackNo, algorithm)
		gpxDoc.MovementStats.MovingData.SetValuesFromMovementData(&gpxTrack.MovementStats.MovingData, trackNo, algorithm)
		gpxDoc.MovementStats.StoppedData.SetValuesFromMovementData(&gpxTrack.MovementStats.StoppedData, trackNo, algorithm)
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
			// TODO:
			//r.Links = route.Links
			if route.Number.Null() {
				r.Number = routeNo
			} else {
				r.Number = route.Number.Value()
			}
			r.Type = route.Type
			// TODO:
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
	time, err := parseGPXTime(original.Timestamp)
	if err == nil {
		result.Timestamp.SetTime(time)
	}
	result.MagneticVariation = original.MagVar
	result.GeoidHeight = original.GeoIdHeight
	result.Name = original.Name
	result.Comment = original.Cmt
	result.Description = original.Desc
	result.Source = original.Src
	// TODO:
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
	if original.DGpsID != nil {
		result.DGpsID = *generic.NewNullableInt(*original.DGpsID)
	}
	return result
}

func convertPointToGpx00(original *geo.GPXPoint) *GPX00GpxPoint {
	result := new(GPX00GpxPoint)
	result.Lat = original.Latitude
	result.Lon = original.Longitude
	result.Ele = original.Elevation
	if original.Timestamp.Valid {
		result.Timestamp = formatGPXTime(original.Timestamp.Time)
	}
	result.MagVar = original.MagneticVariation
	result.GeoIdHeight = original.GeoidHeight
	result.Name = original.Name
	result.Cmt = original.Comment
	result.Desc = original.Description
	result.Src = original.Source
	// TODO:
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
	if original.DGpsID.NotNull() {
		value := original.DGpsID.Value()
		result.DGpsID = &value
	}
	return result
}
