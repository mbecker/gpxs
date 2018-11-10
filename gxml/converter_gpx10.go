// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gxml

import (
	"time"

	"github.com/mbecker/gpxs/geo"
)

//defaultCreator contains the original repo path
const defaultCreator = "https://github.com/tkrajina/gpxgo"

// ----------------------------------------------------------------------------------------------------
// Gpx 1.0 Stuff
// ----------------------------------------------------------------------------------------------------

func convertToGpx10Models(gpxDoc *geo.GPX) *GPX10Gpx {
	gpx10Doc := &GPX10Gpx{}

	//gpx10Doc.XMLNs = gpxDoc.XMLNs
	gpx10Doc.XMLNs = "http://www.topografix.com/GPX/1/0"
	gpx10Doc.XmlNsXsi = gpxDoc.XmlNsXsi
	gpx10Doc.XmlSchemaLoc = gpxDoc.XmlSchemaLoc

	gpx10Doc.Version = "1.0"
	if len(gpxDoc.Creator) == 0 {
		gpx10Doc.Creator = defaultCreator
	} else {
		gpx10Doc.Creator = gpxDoc.Creator
	}
	gpx10Doc.Name = gpxDoc.Name
	gpx10Doc.Desc = gpxDoc.Description
	gpx10Doc.Author = gpxDoc.AuthorName
	gpx10Doc.Email = gpxDoc.AuthorEmail

	if len(gpxDoc.AuthorLink) > 0 || len(gpxDoc.AuthorLinkText) > 0 {
		// ToDo
	}

	if len(gpxDoc.Link) > 0 || len(gpxDoc.LinkText) > 0 {
		gpx10Doc.Url = gpxDoc.Link
		gpx10Doc.UrlName = gpxDoc.LinkText
	}

	if !gpxDoc.Timestamp.IsZero() {
		gpx10Doc.Time = formatGPXTime(gpxDoc.Timestamp)
	}

	gpx10Doc.Keywords = gpxDoc.Keywords

	if gpxDoc.Waypoints != nil {
		gpx10Doc.Waypoints = make([]*GPX00GpxPoint, len(gpxDoc.Waypoints))
		for waypointNo, waypoint := range gpxDoc.Waypoints {
			gpx10Doc.Waypoints[waypointNo] = convertPointToGpx00(&waypoint)
		}
	}

	if gpxDoc.Routes != nil {
		gpx10Doc.Routes = make([]*GPX00GpxRte, len(gpxDoc.Routes))
		for routeNo, route := range gpxDoc.Routes {
			r := new(GPX00GpxRte)
			r.Name = route.Name
			r.Cmt = route.Comment
			r.Desc = route.Description
			r.Src = route.Source
			// ToDo
			//r.Links = route.Links
			r.Number.SetValue(route.Number)
			r.Type = route.Type
			// ToDo
			//r.RoutePoints = route.RoutePoints

			gpx10Doc.Routes[routeNo] = r

			if route.Points != nil {
				r.Points = make([]*GPX00GpxPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointToGpx00(&point)
				}
			}
		}
	}

	if gpxDoc.Tracks != nil {
		gpx10Doc.Tracks = make([]*GPX00GpxTrk, len(gpxDoc.Tracks))
		for trackNo, track := range gpxDoc.Tracks {
			gpx10Track := new(GPX00GpxTrk)
			gpx10Track.Name = track.Name
			gpx10Track.Cmt = track.Comment
			gpx10Track.Desc = track.Description
			gpx10Track.Src = track.Source
			gpx10Track.Number.SetValue(track.Number)

			gpx10Track.Type = track.Type

			if track.Segments != nil {
				gpx10Track.Segments = make([]*GPX00GpxTrkSeg, len(track.Segments))
				for segmentNo, segment := range track.Segments {
					gpx10Segment := new(GPX00GpxTrkSeg)
					if segment.Points != nil {
						gpx10Segment.Points = make([]*GPX00GpxPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpx10Point := convertPointToGpx00(&point)
							// ToDo
							//gpx10Point.Speed = point.Speed
							//gpx10Point.Speed = point.Speed
							gpx10Segment.Points[pointNo] = gpx10Point
						}
					}
					gpx10Track.Segments[segmentNo] = gpx10Segment
				}
			}
			gpx10Doc.Tracks[trackNo] = gpx10Track
		}
	}

	return gpx10Doc
}

func convertFromGpx10Models(gpx10Doc *GPX10Gpx, algorithm geo.Algorithm) *geo.GPX {
	gpxDoc := new(geo.GPX)

	gpxDoc.XMLNs = gpx10Doc.XMLNs
	gpxDoc.XmlNsXsi = gpx10Doc.XmlNsXsi
	gpxDoc.XmlSchemaLoc = gpx10Doc.XmlSchemaLoc

	gpxDoc.Creator = gpx10Doc.Creator
	gpxDoc.Version = gpx10Doc.Version
	gpxDoc.Name = gpx10Doc.Name
	gpxDoc.Description = gpx10Doc.Desc
	gpxDoc.AuthorName = gpx10Doc.Author
	gpxDoc.AuthorEmail = gpx10Doc.Email

	if len(gpx10Doc.Url) > 0 || len(gpx10Doc.UrlName) > 0 {
		gpxDoc.Link = gpx10Doc.Url
		gpxDoc.LinkText = gpx10Doc.UrlName
	}

	// ToDo: If no Time is given then use the first point's time
	if len(gpx10Doc.Time) > 0 {
		gpxDoc.Timestamp, _ = parseGPXTime(gpx10Doc.Time)
	} else {
		gpxDoc.Timestamp = new(time.Time)
	}

	gpxDoc.Keywords = gpx10Doc.Keywords

	Set00GPX00DocWaypoint(gpxDoc, gpx10Doc.Waypoints)
	Set00GPX00DocRoutes(gpxDoc, gpx10Doc.Routes)

	// 1.) Copy if exists GPX00Tracks (gpx10Doc.Tracks) to gpxDoc.Tracks; 2.) If the gpxDoc.Name is empty the assign the track name (FIFO)
	Converter00GPX00DocTracks(gpxDoc, gpx10Doc.Tracks, algorithm)

	return gpxDoc
}
