// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gxml

import (
	"encoding/xml"

	"github.com/mbecker/gpxs/geo"
)

/*

The GPX XML hierarchy:

gpx
    - attr: version (xsd:string) required
    - attr: creator (xsd:string) required
    name
    desc
    author
    email
    url
    urlname
    time
    keywords
    bounds
    wpt
        - attr: lat (gpx:latitudeType) required
        - attr: lon (gpx:longitudeType) required
        ele
        time
        magvar
        geoidheight
        name
        cmt
        desc
        src
        url
        urlname
        sym
        type
        fix
        sat
        hdop
        vdop
        pdop
        ageofdgpsdata
        dgpsid
    rte
        name
        cmt
        desc
        src
        url
        urlname
        number
        rtept
            - attr: lat (gpx:latitudeType) required
            - attr: lon (gpx:longitudeType) required
            ele
            time
            magvar
            geoidheight
            name
            cmt
            desc
            src
            url
            urlname
            sym
            type
            fix
            sat
            hdop
            vdop
            pdop
            ageofdgpsdata
            dgpsid
    trk
        name
        cmt
        desc
        src
        url
        urlname
        number
        trkseg
            trkpt
                - attr: lat (gpx:latitudeType) required
                - attr: lon (gpx:longitudeType) required
                ele
                time
                course
                speed
                magvar
                geoidheight
                name
                cmt
                desc
                src
                url
                urlname
                sym
                type
                fix
                sat
                hdop
                vdop
                pdop
                ageofdgpsdata
                dgpsid
*/

type GPX10Gpx struct {
	XMLName      xml.Name `xml:"gpx"`
	XMLNs        string   `xml:"xmlns,attr,omitempty"`
	XMLNsXsi     string   `xml:"xmlns:xsi,attr,omitempty"`
	XMLSchemaLoc string   `xml:"xsi:schemaLocation,attr,omitempty"`

	Version   string           `xml:"version,attr"`
	Creator   string           `xml:"creator,attr"`
	Name      string           `xml:"name,omitempty"`
	Desc      string           `xml:"desc,omitempty"`
	Author    string           `xml:"author,omitempty"`
	Email     string           `xml:"email,omitempty"`
	Url       string           `xml:"url,omitempty"`
	UrlName   string           `xml:"urlname,omitempty"`
	Time      string           `xml:"time,omitempty"`
	Keywords  string           `xml:"keywords,omitempty"`
	Bounds    *geo.GpxBounds   `xml:"bounds"`
	Waypoints []*GPX00GpxPoint `xml:"wpt"`
	Routes    []*GPX00GpxRte   `xml:"rte"`
	Tracks    []*GPX00GpxTrk   `xml:"trk"`
}
