## Snapshot intervals and frame durations

A good default is  0.04167

For a day (24 hours):
- Snapshot interval: 1 minute
- Duration per frame: 0.04167 (1/24 second)
- Total frames: 1440
- Final video length: ~60 seconds

For a week:
- Snapshot interval: 5 minutes
- Duration per frame: 0.04167
- Total frames: 2016
- Final video length: ~84 seconds

For a month:
- Snapshot interval: 15 minutes
- Duration per frame: 0.04167
- Total frames: 2880
- Final video length: ~120 seconds

For 6 months:
- Snapshot interval: 1 hour
- Duration per frame: 0.04167
- Total frames: 4320
- Final video length: ~180 seconds

For a year:
- Snapshot interval: 2 hours
- Duration per frame: 0.04167
- Total frames: 4380
- Final video length: ~183 seconds

Gives smooth playback at standard 24fps. If you want to adjust the final video
length, you can modify either the capture interval or the frame duration.

0.08333 (1/12 second)
- Creates a slightly slower, more contemplative feel

0.0333 (1/30 second)
- Slightly faster, more dynamic feel

0.0208 (1/48 second)
- Creates very smooth motion

For capture intervals, some alternative useful values:
- 30 seconds: Good for fast-changing scenes like sunset/sunrise
- 2 minutes: Nice for cloud movements
- 10 minutes: Works well for construction sites
- 3 hours: Good for seasonal changes
- 4 hours: Nice for garden/plant growth
- 12 hours: Captures day/night cycles effectively
