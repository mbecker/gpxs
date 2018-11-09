# gpxs - gpxgo on steriods

## Tests - ToDo
- Test that GPX Distance - Sum of all Track Distance - Sum of all Tracks Segments Distance - Sm of all Tracks Segments MovingData Distance are equal
- Test that GPX Duration - Sum of all Track Duration - Sum of all Tracks Segments Duration - Sm of all Tracks Segments MovingData Duration are equal

## Bechmark: 1046 gpx files
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


