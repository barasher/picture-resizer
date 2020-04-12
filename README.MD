# picture-resizer

[![GoDoc](https://godoc.org/github.com/barasher/picture-resizer?status.svg)](https://godoc.org/github.com/barasher/picture-resizer)

**picture-resizer** is a tools that resizes a batch of picture files.

It is built on top of [ImageMagick](https://imagemagick.org/index.php).

## Binary version

### Requirements

The binary version requires ImageMagick's MagickWand C API.

Installation instructions :
- On Ubuntu / Debian : `sudo apt-get install libmagickwand-dev`

### Release

Every releases are available on **picture-resizer**'s [Github page](https://github.com/barasher/picture-resizer/releases).

### From source

TODO

### Execution

```shell script
$ ./picture-resizer -h
Picture resizer

Usage:
  picture-resizer [flags]

Flags:
  -b, --blur float            blur (> 1 is blurry, < 1 is sharp) (default 1)
  -c, --compression uint      compression quality (1-100) (default 95)
  -H, --height uint           height (pixels) (default 480)
  -h, --help                  help for picture-resizer
  -i, --input string          input folder
  -l, --loggingLevel string   logging level (default "info")
  -o, --output string         output folder
  -t, --threadCount uint      thread count (default 2)
  -w, --width uint            width (pixels) (default 640)
```

The only required parameters are `-i` and `-o`.

## Docker version

### From Docker Hub

### From Dockerfile

### Execution

## Logs

```shell script
$ ./picture-resizer -i /tmp/input/ -o /tmp/output -l debug
```

```json    
{"level":"info","threadId":0,"time":"2020-04-12T16:11:32+02:00","message":"Converting /tmp/input/pic1.jpg..."}
{"level":"info","threadId":1,"time":"2020-04-12T16:11:32+02:00","message":"Converting /tmp/input/pic2.jpg..."}
{"level":"debug","threadId":0,"time":"2020-04-12T16:11:32+02:00","message":"Output for /tmp/input/pic1.jpg: /tmp/output/1afd9136b5ea6e72c48df7c965a1d24b__pic1.jpg"}
{"level":"debug","threadId":1,"time":"2020-04-12T16:11:33+02:00","message":"Output for /tmp/input/pic2.jpg: /tmp/output/b165ee33473dc2a1902dd5ced13cf97d__pic2.jpg"}
```