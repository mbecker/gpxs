// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gxml

import (
	"encoding/xml"
)

/*

The GPX XML hierarchy:

gpx (gpxType)
    - attr: version (xsd:string) None
    - attr: creator (xsd:string) None
    metadata (metadataType)
        name (xsd:string)
        desc (xsd:string)
        author (personType)
            name (xsd:string)
            email (emailType)
                - attr: id (xsd:string) None
                - attr: domain (xsd:string) None
            link (linkType)
                - attr: href (xsd:anyURI) None
                text (xsd:string)
                type (xsd:string)
        copyright (copyrightType)
            - attr: author (xsd:string) None
            year (xsd:gYear)
            license (xsd:anyURI)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        time (xsd:dateTime)
        keywords (xsd:string)
        bounds (boundsType)
            - attr: minlat (latitudeType) None
            - attr: minlon (longitudeType) None
            - attr: maxlat (latitudeType) None
            - attr: maxlon (longitudeType) None
        extensions (extensionsType)
    wpt (wptType)
        - attr: lat (latitudeType) None
        - attr: lon (longitudeType) None
        ele (xsd:decimal)
        time (xsd:dateTime)
        magvar (degreesType)
        geoidheight (xsd:decimal)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        sym (xsd:string)
        type (xsd:string)
        fix (fixType)
        sat (xsd:nonNegativeInteger)
        hdop (xsd:decimal)
        vdop (xsd:decimal)
        pdop (xsd:decimal)
        ageofdgpsdata (xsd:decimal)
        dgpsid (dgpsStationType)
        extensions (extensionsType)
    rte (rteType)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        number (xsd:nonNegativeInteger)
        type (xsd:string)
        extensions (extensionsType)
        rtept (wptType)
            - attr: lat (latitudeType) None
            - attr: lon (longitudeType) None
            ele (xsd:decimal)
            time (xsd:dateTime)
            magvar (degreesType)
            geoidheight (xsd:decimal)
            name (xsd:string)
            cmt (xsd:string)
            desc (xsd:string)
            src (xsd:string)
            link (linkType)
                - attr: href (xsd:anyURI) None
                text (xsd:string)
                type (xsd:string)
            sym (xsd:string)
            type (xsd:string)
            fix (fixType)
            sat (xsd:nonNegativeInteger)
            hdop (xsd:decimal)
            vdop (xsd:decimal)
            pdop (xsd:decimal)
            ageofdgpsdata (xsd:decimal)
            dgpsid (dgpsStationType)
            extensions (extensionsType)
    trk (trkType)
        name (xsd:string)
        cmt (xsd:string)
        desc (xsd:string)
        src (xsd:string)
        link (linkType)
            - attr: href (xsd:anyURI) None
            text (xsd:string)
            type (xsd:string)
        number (xsd:nonNegativeInteger)
        type (xsd:string)
        extensions (extensionsType)
        trkseg (trksegType)
            trkpt (wptType)
                - attr: lat (latitudeType) None
                - attr: lon (longitudeType) None
                ele (xsd:decimal)
                time (xsd:dateTime)
                magvar (degreesType)
                geoidheight (xsd:decimal)
                name (xsd:string)
                cmt (xsd:string)
                desc (xsd:string)
                src (xsd:string)
                link (linkType)
                    - attr: href (xsd:anyURI) None
                    text (xsd:string)
                    type (xsd:string)
                sym (xsd:string)
                type (xsd:string)
                fix (fixType)
                sat (xsd:nonNegativeInteger)
                hdop (xsd:decimal)
                vdop (xsd:decimal)
                pdop (xsd:decimal)
                ageofdgpsdata (xsd:decimal)
                dgpsid (dgpsStationType)
                extensions (extensionsType)
            extensions (extensionsType)
    extensions (extensionsType)
*/

type GPX11Gpx struct {
	XMLName      xml.Name `xml:"gpx"`
	XMLNs        string   `xml:"xmlns,attr,omitempty"`
	XmlNsXsi     string   `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string   `xml:"xsi:schemaLocation,attr,omitempty"`

	Version     string         `xml:"version,attr"`
	Creator     string         `xml:"creator,attr"`
	Name        string         `xml:"metadata>name,omitempty"`
	Desc        string         `xml:"metadata>desc,omitempty"`
	AuthorName  string         `xml:"metadata>author>name,omitempty"`
	AuthorEmail *GPX11GpxEmail `xml:"metadata>author>email,omitempty"`
	/**
	 * TODO:
	 * - [x] Can be me more tane one link for AuthorLink?
	 */
	AuthorLink *GPX00GpxLink       `xml:"metadata>author>link,omitempty"`
	Copyright  *GPX11GpxCopyright  `xml:"metadata>copyright,omitempty"`
	Link       *GPX00GpxLink       `xml:"metadata>link,omitempty"`
	Timestamp  string              `xml:"metadata>time,omitempty"`
	Keywords   string              `xml:"metadata>keywords,omitempty"`
	Bounds     *GPX11GpxBounds     `xml:"bounds"`
	Extensions *GPX11GpxExtensions `xml:"extensions"`
	Waypoints  []*GPX00GpxPoint    `xml:"wpt"`
	Routes     []*GPX00GpxRte      `xml:"rte"`
	Tracks     []*GPX00GpxTrk      `xml:"trk"`
}

type GPX11GpxBounds struct {
	//XMLName xml.Name `xml:"bounds"`
	MinLat float64 `xml:"minlat,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

type GPX11GpxCopyright struct {
	XMLName xml.Name `xml:"copyright"`
	Author  string   `xml:"author,attr"`
	Year    string   `xml:"year,omitempty"`
	License string   `xml:"license,omitempty"`
}

//type GPX11GpxAuthor struct {
//	Name  string        `xml:"name,omitempty"`
//	Email string        `xml:"email,omitempty"`
//	Link  *GPX10GpxLink `xml:"link"`
//}

type GPX11GpxEmail struct {
	Id     string `xml:"id,attr"`
	Domain string `xml:"domain,attr"`
}

//type GPX11GpxMetadata struct {
//	XMLName xml.Name        `xml:"metadata"`
//	Name    string          `xml:"name,omitempty"`
//	Desc    string          `xml:"desc,omitempty"`
//	Author  *GPX11GpxAuthor `xml:"author,omitempty"`
//	//	Copyright *GpxCopyright `xml:"copyright,omitempty"`
//	//	Links     []GpxLink     `xml:"link"`
//	Timestamp string `xml:"time,omitempty"`
//	Keywords  string `xml:"keywords,omitempty"`
//	//	Bounds    *GpxBounds    `xml:"bounds"`
//}

type GPX11GpxExtensions struct {
	Bytes []byte `xml:",innerxml"`
}

/**
 * Common struct fields for all points
 */
// type GPX10GpxPoint struct {
// 	Lat float64 `xml:"lat,attr"`
// 	Lon float64 `xml:"lon,attr"`
// 	// Position info
// 	Ele         generic.NullableFloat64 `xml:"ele,omitempty"`
// 	Timestamp   string                  `xml:"time,omitempty"`
// 	MagVar      string                  `xml:"magvar,omitempty"`
// 	GeoIdHeight string                  `xml:"geoidheight,omitempty"`
// 	// Description info
// 	Name  string         `xml:"name,omitempty"`
// 	Cmt   string         `xml:"cmt,omitempty"`
// 	Desc  string         `xml:"desc,omitempty"`
// 	Src   string         `xml:"src,omitempty"`
// 	Links []GPX10GpxLink `xml:"link"`
// 	Sym   string         `xml:"sym,omitempty"`
// 	Type  string         `xml:"type,omitempty"`
// 	// Accuracy info
// 	Fix           string   `xml:"fix,omitempty"`
// 	Sat           *int     `xml:"sat,omitempty"`
// 	Hdop          *float64 `xml:"hdop,omitempty"`
// 	Vdop          *float64 `xml:"vdop,omitempty"`
// 	Pdop          *float64 `xml:"pdop,omitempty"`
// 	AgeOfDGpsData *float64 `xml:"ageofdgpsdata,omitempty"`
// 	DGpsId        *int     `xml:"dgpsid,omitempty"`
// }
