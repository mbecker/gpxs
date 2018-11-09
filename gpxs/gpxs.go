package gpxs

import (
	"github.com/mbecker/gpxs/geo"
	gxml "github.com/mbecker/gpxs/gxml"
)

/* ToDo: Singleton - Does it make sense?
type options struct {
	algorithm geo.Algorithm
}

// Init singleton
var option *options
var once sync.Once

func Option() *options {
	once.Do(func() {
		option = &options{}
	})
	return option
}

func (opt *options) SetAlgorithm(algorithm geo.Algorithm) {
	opt.algorithm = algorithm
}
*/

type ChannelStruct struct {
	GPX   *geo.GPX
	Error error
}

func ParseFileChannel(fileName string, algorithm geo.Algorithm, c chan ChannelStruct) {
	gpx, err := gxml.ParseFile(fileName, algorithm)
	if err != nil {
		channelStruct := ChannelStruct{
			GPX:   nil,
			Error: err,
		}
		c <- channelStruct
	} else {
		channelStruct := ChannelStruct{
			GPX:   gpx,
			Error: nil,
		}
		c <- channelStruct
	}
}

//ParseFile parses a gpx file and returns a GPX object
func ParseFile(fileName string, algorithm geo.Algorithm) (*geo.GPX, error) {
	return gxml.ParseFile(fileName, algorithm)
}

//ParseBytes parses GPX from bytes
func ParseBytes(bytes []byte, algorithm geo.Algorithm) (*geo.GPX, error) {
	return gxml.ParseBytes(bytes, algorithm)
}

//ParseString parses GPX from string
func ParseString(str string, algorithm geo.Algorithm) (*geo.GPX, error) {
	return gxml.ParseBytes([]byte(str), algorithm)
}
