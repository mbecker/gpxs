# gpxs - gpxgo on steriods

> All credits to gpx gpxgo - https://github.com/tkrajina/gpxgo/issues

# Motivation
I'm using Strava quite heavily (3-5 activities/week) for most of my sport activities like running, cycling, hiking, etc. (see https://www.strava.com/athletes/8844168). At the moment I'm using the Strava iOS App.
Unfortunatley I started with Runkeeper years ago and switched for a time to Garmin and other tools. Once I exported (more or less all) my tracks from all third parties and uploaded the gpx files to Strava. Strava analyzed all the files for me and added the activities to my statistics. Nice!

I want to have all my activites on my local desktop to analyze the gpx files and create nice stats and graphics. For that I downloaded all my gpx files from Strava and started to look how to parse the gpx files. That's how I discovered __gpxgo__.

(How to download all activities see: https://support.strava.com/hc/en-us/articles/216918437-Exporting-your-Data-and-Bulk-Export #Section: Bulk Export at the the end of the page)

The fantastic library __gpxgo__ (https://github.com/tkrajina/gpxgo/issues) by tkrajina parses gpx files and returns information of the gpx/tracks/segements/points like distance, duration and movingtime/distance + stoppedtime/distance.

Unfortunatley the gpx information like distance (movingdistance), duration (movingtime) does not match with the Strava information.

My assumption is that Strava does a quite good in analyzing gpx data and should be the baseline to match. I looked into the libraray and noticed for example that the distance calculation is quite fast but less accurate than other GPS distance clulcation formulas like Haversine and Vincenty. I started to implement these formulas in the libraray and quickly noticed that for each part of the gpx (gpx->track->segment->point) the data is re-calculated new. Additionaly the points which are used for the caluclation for "Moving"Time/Distance are based on the assumption that the speed between two points must be greater than 1.0m/s. That results in a big difference of "Moving"Time/Distance of Strava and the library.

**So I've started to update te libraray as follows:**

- Change the method how the data is aggregated (from point->segments->tracks->gpx)
- Refactoring of the code into multiple packages; refactoring the XML strucst/converter since the gpx v1.0 and v1.1 isn't so different
- Use track name as gpx name if it's not already set; Use first point timestamp for track/gpx timestamp if it's not already set
- User provides the the caluclation of distance, duration, speed, pace, etc. + method to normalize the points
- Added new GPS distance cluclation
- Use standard deviation to determine which points should be used for "Moving"Time/Distance

**How is the data now agregated?**

At first the point data like duration, distane between the previous poind and the actual point is caluclated (with the user provided methods). The the data is then aggregated to the next level

```
Point(distance and duration previous and actual point)
All Points(sum of distance and duration) -> Segment
All segments(sum of distance and duration) -> Track
All tracks(sum of distance and duration) -> Gpx
```

**Refactoring**

For my personal style to work with multiple structs, interface and funcs I splitted the code into multiple packages. That makes it easier for me to work with and to understand the code.
The structs which define the XML structure of gpx v1.0 and v1.0 are quite similar. So I tried to use structs as a baseline (call it v00) and then only implement structs for the specifc XML elements which differs. With that the conversion from XML to GPX (and vice versa) is only implemented once (and not twice).

**GPX name and timestamp from track and first point**

In the gpx files from Strava (and the other tird parties I've used) no name and sometimes no timestamp is given for gpx at the top level. If that's the case then __gpxs__ uses the name of the track(s) and the timestamp of the first point. The name and timestamp is taken by FIFO (first in first out); that means basically that first name/timestamp is used which exists.

**User method to calculate distance, duration, speed, pace + method to normalize the points**

The idea is that the uses who uses this libraray provides own methods to calculate the points data. That means the users could use a GPS distance calculatio not yet provided by the library. Additionaly you could use your own standardize method o use only the points you want to use for your calculation.
Methods implemented:
- "Standard Method" by __gpxgo__: Basic but fast calculation of distance between two points
- Vincenty: More accurate methods
- Normalize method (1): Threshold by __gpxgo__; use only points above speed threshold of 1m/s
- Normalize method (2): Standard deviation (see below)

**New GPS distance cluclation**

Added a method to use GPS formula Vincenty. Also added the calculation methods from __gpxgo__ as the "Standard Algorithm".

**Standard deviation for "Moving"Time/Distance?**

A point has the following data: lat/long (where) and timestamp (when). With the information where and when of each point you could calculate the distance (how long in meters?) and duration (how long in seconds?) between the previous and actual point.

It's quite normal that you stop during an activity (waiting on the street, doing a longer pause, etc) and then the distance/duration between the point before and after the pause could be quite long. Let's say your are doing a cycling tour to a castle; at the castle you stop (Point p0(position0, time0)) and after 1hour you are cycling back (start of the return Point p1(position0+100m, timee0+3600seconds)). The speed would be 

```
P1(position0+100m)-P0(position0) / P1(time0+36000) - P0(time0) -> 100m / 3600sec -> ~0,0278 m/s = 0,10008 km/h
```

Obvious the data between these two point p0 and p1 shouldn't be used for the "Moving"Time/Distance since p1 is a new startig point. The point after P1 should then be used for the "Moving"Time/Distance.

The idea of gpxgo is that the speed between two points must be grater than 1m/s; then the points are used for the caluclation. The approach is fair enough since it's quite fast. But for me the different to Strava is too much (for example the duration for a cycling tour with gpxgo is 3h20min and with Stava it's 2h45min which is the "real" movingtime from my trip).

The standard deviation (https://en.wikipedia.org/wiki/Standard_deviation) uses all points data of the gpx track and "quantifies the amount of variation or dispersion of the set of data values". In my words: You define a weightehd average of all points duration and then define that only the points should be used which are close enough to that duration average (like 90% or 95%).
With an standard deviation Sigma σ of 1.644854 (~95%) I'm quite close to the Strava information (only some seconds in each of my ~1050 gpx files.

So good enough for me!

**Performance**
See benchmarks below.

My assumption was that the "new" libraray gpxs is faster than __gpxgo__ because of the new aggregation method. For a benchmark of Vincenty I added the method to __gpxgo__. __gpxgo__ does not use the standard deviation; so __gpxgs__ is instrumented to use the same "normalization method" as __gpxgo__.


__Basic distance cluclation; no standard deviation; default speed threshold of 1m/sec (gpx files: 1046)__
gpxs (3) 17.2373831s vs. gpxgo (1) 19.1264728s

__Vincenty distance calculation; no standard deviation; default speed threshold of 1m/sec (gpx files: 1046)__
gpxs (2) 23.014476s vs. gpxgo (2) 25.3617024s

__Vincenty distance caluclation; standard deviation of 1.644854__
gpxs (1) 21.5685324s vs. pxgo (not implemented)

**Conclusion**

The reuslts of gpxs (distance, "moving"time/distance, "stopped"time/distance) are for **me more accurate**. "More accurate" **means for me** the calculated infomation is closer to Stava which is **my baseline**.

The "more accurate" is achieved by using a new GPS distance calculation (Vincenty formula) and the standard deviation to nromalize the data set (which points should be used for the "moving"time/distance calculation).

Suprisingly the new aggregation is not so faster as I thought ("only" 2seconds); using the standard deviation it's faster (4sec).

Additionaly the new approach to use custom methods to normalize the data sete and caluclate the information provides more flexibility.


# ToDo

- Test that GPX Distance - Sum of all Track Distance - Sum of all Tracks Segments Distance - Sm of all Tracks Segments MovingData Distance are equal
- Test that GPX Duration - Sum of all Track Duration - Sum of all Tracks Segments Duration - Sm of all Tracks Segments MovingData Duration are equal
- Write test files
- Write benchmark fies

# Bechmark: Different gps distance calculations

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
