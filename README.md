## Configuration

Example configuration.

```yaml
cameras:
  - name: "Riksgransen"         # Name of the camera, if using spaces in name it will be converted: Hello world -> helloWorld
    snapshotUrl: "https://.."   # URL
    auth:                       # If snapshotUrl need authentication
      type: "basic"             # Can be basic or bearer
      username: "user"          # Username (basic auth)
      password: "pass"          # Password (basic auth)
    interval: "*/10 * * * *"    # Snapshot interval cron expression
    timelapseInterval: "* 24,12 * * * *" # Timelapse generaton cron expression interval
    delete: true                # Delete snapshot images after timelapse generation
    frameDuration: 0.041667     # Frame duration for each snapshot
    ffmpeg_template: "ffmpeg ... -i {{.ListPath}} ... -y {{.OutputPath}}" # ffmpeg command used for timelapse generation.

# Where to write snapshots and timelapses
outputDir: "/tmp/timelapser"
# Defaults used if not set per camera
interval: "*/5 * * * *"
timelapseInterval: "* 24,12 * * * *"
frameDuration: 0.041667
ffmpeg_template: "ffmpeg -f concat -safe 0 -i {{.ListPath}} -vf fps=24,format=yuv420p -c:v libx264 -preset medium -crf 23 -movflags +faststart -y {{.OutputPath}}"
```

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


## How to integrate timelapser with Home Assistant

Integration with Home Assistant is easy by using the the Home Assistant Local
Media feautre.

0 8 * * * mv /media/timelapse/*.mp4 /usr/share/hassio/media

After this the timelapses are available every day at 8 AM UTC to be viewed using
the Local Media browser, or even better using the Gallery card like this:

```yaml
type: 'custom:gallery-card'
entities:
  - path: 'media-source://media_source/media/'
    recursive: true
menu_alignment: Hidden
file_name_format: '*.mp4'
```
