package gxml

import (
	"strings"
	"time"

	"github.com/mbecker/gpxs/geo"
)

// ----------------------------------------------------------------------------------------------------
// Gpx 1.1 Stuff
// ----------------------------------------------------------------------------------------------------

func convertToGpx11Models(gpxDoc *geo.GPX) *GPX11Gpx {
	gpx11Doc := &GPX11Gpx{}

	gpx11Doc.Version = "1.1"

	gpx11Doc.XMLNs = "http://www.topografix.com/GPX/1/1"
	gpx11Doc.XMLNsXsi = gpxDoc.XMLNsXsi
	gpx11Doc.XMLSchemaLoc = gpxDoc.XMLSchemaLoc

	if len(gpxDoc.Creator) == 0 {
		gpx11Doc.Creator = defaultCreator
	} else {
		gpx11Doc.Creator = gpxDoc.Creator
	}
	gpx11Doc.Name = gpxDoc.Name
	gpx11Doc.Desc = gpxDoc.Description
	gpx11Doc.AuthorName = gpxDoc.AuthorName

	if len(gpxDoc.AuthorEmail) > 0 {
		parts := strings.Split(gpxDoc.AuthorEmail, "@")
		if len(parts) == 1 {
			gpx11Doc.AuthorEmail = new(GPX11GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
		} else if len(parts) > 1 {
			gpx11Doc.AuthorEmail = new(GPX11GpxEmail)
			gpx11Doc.AuthorEmail.Id = parts[0]
			gpx11Doc.AuthorEmail.Domain = parts[1]
		}
	}

	if len(gpxDoc.AuthorLink) > 0 || len(gpxDoc.AuthorLinkText) > 0 || len(gpxDoc.AuthorLinkType) > 0 {
		// ToDO: Why not GPX10GpxLink? The type does not match?!
		gpx11Doc.AuthorLink = new(GPX00GpxLink)
		gpx11Doc.AuthorLink.Href = gpxDoc.AuthorLink
		gpx11Doc.AuthorLink.Text = gpxDoc.AuthorLinkText
		gpx11Doc.AuthorLink.Type = gpxDoc.AuthorLinkType
	}

	if len(gpxDoc.Copyright) > 0 || len(gpxDoc.CopyrightYear) > 0 || len(gpxDoc.CopyrightLicense) > 0 {
		gpx11Doc.Copyright = new(GPX11GpxCopyright)
		gpx11Doc.Copyright.Author = gpxDoc.Copyright
		gpx11Doc.Copyright.Year = gpxDoc.CopyrightYear
		gpx11Doc.Copyright.License = gpxDoc.CopyrightLicense
	}

	if len(gpxDoc.Link) > 0 || len(gpxDoc.LinkText) > 0 || len(gpxDoc.LinkType) > 0 {
		gpx11Doc.Link = new(GPX00GpxLink)
		gpx11Doc.Link.Href = gpxDoc.Link
		gpx11Doc.Link.Text = gpxDoc.LinkText
		gpx11Doc.Link.Type = gpxDoc.LinkType
	}

	if !gpxDoc.Timestamp.IsZero() {
		gpx11Doc.Timestamp = formatGPXTime(gpxDoc.Timestamp)
	}

	gpx11Doc.Keywords = gpxDoc.Keywords

	if gpxDoc.Waypoints != nil {
		gpx11Doc.Waypoints = make([]*GPX00GpxPoint, len(gpxDoc.Waypoints))
		for waypointNo, waypoint := range gpxDoc.Waypoints {
			gpx11Doc.Waypoints[waypointNo] = convertPointToGpx00(&waypoint)
		}
	}

	if gpxDoc.Routes != nil {
		gpx11Doc.Routes = make([]*GPX00GpxRte, len(gpxDoc.Routes))
		for routeNo, route := range gpxDoc.Routes {
			r := new(GPX00GpxRte)
			r.Name = route.Name
			r.Cmt = route.Comment
			r.Desc = route.Description
			r.Src = route.Source
			// TODO:
			//r.Links = route.Links
			r.Number.SetValue(route.Number)
			r.Type = route.Type
			// TODO:
			//r.RoutePoints = route.RoutePoints

			gpx11Doc.Routes[routeNo] = r

			if route.Points != nil {
				r.Points = make([]*GPX00GpxPoint, len(route.Points))
				for pointNo, point := range route.Points {
					r.Points[pointNo] = convertPointToGpx00(&point)
				}
			}
		}
	}

	if gpxDoc.Tracks != nil {
		gpx11Doc.Tracks = make([]*GPX00GpxTrk, len(gpxDoc.Tracks))
		for trackNo, track := range gpxDoc.Tracks {
			gpx11Track := new(GPX00GpxTrk)
			gpx11Track.Name = track.Name
			gpx11Track.Cmt = track.Comment
			gpx11Track.Desc = track.Description
			gpx11Track.Src = track.Source
			gpx11Track.Number.SetValue(track.Number)
			gpx11Track.Type = track.Type

			if track.Segments != nil {
				gpx11Track.Segments = make([]*GPX00GpxTrkSeg, len(track.Segments))
				for segmentNo, segment := range track.Segments {
					gpx11Segment := new(GPX00GpxTrkSeg)
					if segment.Points != nil {
						gpx11Segment.Points = make([]*GPX00GpxPoint, len(segment.Points))
						for pointNo, point := range segment.Points {
							gpx11Segment.Points[pointNo] = convertPointToGpx00(&point)
						}
					}
					gpx11Track.Segments[segmentNo] = gpx11Segment
				}
			}
			gpx11Doc.Tracks[trackNo] = gpx11Track
		}
	}

	return gpx11Doc
}

func convertFromGpx11Models(gpx11Doc *GPX11Gpx, algorithm geo.Algorithm) *geo.GPX {

	gpxDoc := new(geo.GPX)

	gpxDoc.XMLNs = gpx11Doc.XMLNs
	gpxDoc.XMLNsXsi = gpx11Doc.XMLNsXsi
	gpxDoc.XMLSchemaLoc = gpx11Doc.XMLSchemaLoc

	gpxDoc.Creator = gpx11Doc.Creator
	gpxDoc.Version = gpx11Doc.Version
	gpxDoc.Name = gpx11Doc.Name
	gpxDoc.Description = gpx11Doc.Desc
	gpxDoc.AuthorName = gpx11Doc.AuthorName

	if gpx11Doc.AuthorEmail != nil {
		gpxDoc.AuthorEmail = gpx11Doc.AuthorEmail.Id + "@" + gpx11Doc.AuthorEmail.Domain
	}
	if gpx11Doc.AuthorLink != nil {
		gpxDoc.AuthorLink = gpx11Doc.AuthorLink.Href
		gpxDoc.AuthorLinkText = gpx11Doc.AuthorLink.Text
		gpxDoc.AuthorLinkType = gpx11Doc.AuthorLink.Type
	}

	/* TODO:
	if gpx11Doc.Extensions != nil {
		gpxDoc.Extensions = &gpx11Doc.Extensions.Bytes
	}
	*/

	// TODO: If no Time is given then use the first point's time
	if len(gpx11Doc.Timestamp) > 0 {
		gpxDoc.Timestamp, _ = parseGPXTime(gpx11Doc.Timestamp)
	} else {
		gpxDoc.Timestamp = new(time.Time)
	}

	if gpx11Doc.Copyright != nil {
		gpxDoc.Copyright = gpx11Doc.Copyright.Author
		gpxDoc.CopyrightYear = gpx11Doc.Copyright.Year
		gpxDoc.CopyrightLicense = gpx11Doc.Copyright.License
	}

	if gpx11Doc.Link != nil {
		gpxDoc.Link = gpx11Doc.Link.Href
		gpxDoc.LinkText = gpx11Doc.Link.Text
		gpxDoc.LinkType = gpx11Doc.Link.Type
	}

	gpxDoc.Keywords = gpx11Doc.Keywords

	Set00GPX00DocWaypoint(gpxDoc, gpx11Doc.Waypoints)
	Set00GPX00DocRoutes(gpxDoc, gpx11Doc.Routes)

	// 1.) Copy if exists GPX00Tracks (gpx11Doc.Tracks) to gpxDoc.Tracks; 2.) If the gpxDoc.Name is empty the assign the track name (FIFO)
	Converter00GPX00DocTracks(gpxDoc, gpx11Doc.Tracks, algorithm)

	return gpxDoc
}
