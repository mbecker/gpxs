# gpxs - gpxgo on steriods

## Tests - ToDo
- Test that GPX Distance - Sum of all Track Distance - Sum of all Tracks Segments Distance - Sm of all Tracks Segments MovingData Distance are equal
- Test that GPX Duration - Sum of all Track Duration - Sum of all Tracks Segments Duration - Sm of all Tracks Segments MovingData Duration are equal

## Bechmark: 1046 gpx files
**gpxs**
```
--- GPX Data ---
Distance: 11534.046010
Duration Time: 1008h26m18s
Moving Distance: 10877.102260 kmMoving Time Time: 893h57m54s
Stopped Distance: 656.943750 km
Stopped Time Time: 114h28m24s
------
2018/11/08 23:18:18 readFiles took 20.3665441s
```
**gpxgo**
```
--- GPX Data ---
Distance: 0.000000
Duration Time: 0s
Moving Distance: 11513.384257 km
Moving Time Time: 962h41m58s
Stopped Distance: 16.449875 km
Stopped Time Time: 45h44m20s
------
2018/11/08 23:19:05 readFiles took 21.4531238s
```


