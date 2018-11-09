# gpxs - gpxgo on steriods

> All credits to gpx gpxgo - https://github.com/tkrajina/gpxgo/issues

# Motivation
I'm using Strava quite heavily (3-5 activities/week) for most my sport activities like running, cycling, hiking, etc. (see https://www.strava.com/athletes/8844168). At the moment I'm using the Strava iOS App.
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
At first the point data like duration, distane between the previous poind and the actual point is caluclated (with the user provided methods). The the data is then aggregated to the next level (point; all points->segment; all segments->track; all tracks->gpx).

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
```
1.) Vincenty - With Standard Deviation
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 10877.102260 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.943750 km
Stopped Time Time: 114h28m24s
--- Tracks ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 10877.102260 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.943750 km
Stopped Time Time: 114h28m24s
--- Segment ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 10877.102260 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.943750 km
Stopped Time Time: 114h28m24s
------
2018/11/09 10:48:26 readFiles took 21.5685324s

2.) Vincenty - Default Stopped Speed Threshold (Speed between points < 1m/s)
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 11453.204095 km
Moving Time Time: 931h9m2s
Stopped Distance: 80.841915 km
Stopped Time Time: 77h17m16s
--- Tracks ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 11453.204095 km
Moving Time Time: 931h9m2s
Stopped Distance: 80.841915 km
Stopped Time Time: 77h17m16s
--- Segment ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 11453.204095 km
Moving Time Time: 931h9m2s
Stopped Distance: 80.841915 km
Stopped Time Time: 77h17m16s
------
2018/11/09 10:37:27 readFiles took 23.014476s

3.) Standard Algorithm (base distance calculation); 2D calculation (no Haversine); no standard deviation
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11506.202686
Duration Time: 1008h26m18s
Moving Distance: 11425.275593 km
Moving Time Time: 931h5m5s
Stopped Distance: 80.927093 km
Stopped Time Time: 77h21m13s
--- Tracks ---
Distance: 11506.202686
Duration Time: 1008h26m18s
Moving Distance: 11425.275593 km
Moving Time Time: 931h5m5s
Stopped Distance: 80.927093 km
Stopped Time Time: 77h21m13s
--- Segment ---
Distance: 11506.202686
Duration Time: 1008h26m18s
Moving Distance: 11425.275593 km
Moving Time Time: 931h5m5s
Stopped Distance: 80.927093 km
Stopped Time Time: 77h21m13s
------
2018/11/09 10:40:51 readFiles took 17.2373831s

4.) Standard Algorithm (base distance calculation); 3D calculation (no Haversine); no standard deviation
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11447.724918 km
Moving Time Time: 931h35m11s
Stopped Distance: 82.109214 km
Stopped Time Time: 76h51m7s
--- Tracks ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11447.724918 km
Moving Time Time: 931h35m11s
Stopped Distance: 82.109214 km
Stopped Time Time: 76h51m7s
--- Segment ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11447.724918 km
Moving Time Time: 931h35m11s
Stopped Distance: 82.109214 km
Stopped Time Time: 76h51m7s
------
2018/11/09 10:38:54 readFiles took 17.0967338s

5.) Standard Algorithm ; Haversine; no standard deviation
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 11445.979743 km
Moving Time Time: 931h8m17s
Stopped Distance: 80.880108 km
Stopped Time Time: 77h18m1s
--- Tracks ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 11445.979743 km
Moving Time Time: 931h8m17s
Stopped Distance: 80.880108 km
Stopped Time Time: 77h18m1s
--- Segment ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 11445.979743 km
Moving Time Time: 931h8m17s
Stopped Distance: 80.880108 km
Stopped Time Time: 77h18m1s
------
2018/11/09 10:43:10 readFiles took 17.4232154s

5.) Standard Algorithm ; Haversine; standard deviation
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 10870.365032 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.494819 km
Stopped Time Time: 114h28m24s
--- Tracks ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 10870.365032 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.494819 km
Stopped Time Time: 114h28m24s
--- Segment ---
Distance: 11526.859851
Duration Time: 1008h26m18s
Moving Distance: 10870.365032 km
Moving Time Time: 893h57m54s
Stopped Distance: 656.494819 km
Stopped Time Time: 114h28m24s
------
2018/11/09 10:45:06 readFiles took 17.6858473s

```
**gpxgo** (Simple distance calculation, no standard deviation for moving data)
```
1.) Length3D() (Base distance clulcation with elevation between points); default stopped threshold
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
--- Tracks ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
--- Segment ---
Distance: 11529.834132
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
------
2018/11/09 09:29:45 readFiles took 19.1264728s

2.) LengthVincenty(); default stopped threshold
--- GPX Files ---
# of Files: 1046
--- GPX ---
Distance: 11.534046
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
--- Tracks ---
Distance: 11.534046
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
--- Segment ---
Distance: 11.534046
Duration Time: 1008h26m18s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
------
2018/11/09 10:51:27 readFiles took 25.3617024s
```


