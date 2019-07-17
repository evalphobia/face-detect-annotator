face-detect-annotator
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/face-detect-annotator?status.svg
[2]: https://godoc.org/github.com/evalphobia/face-detect-annotator
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/face-detect-annotator.svg
[6]: https://github.com/evalphobia/face-detect-annotator/releases/latest
[7]: https://travis-ci.org/evalphobia/face-detect-annotator.svg?branch=master
[8]: https://travis-ci.org/evalphobia/face-detect-annotator
[9]: https://coveralls.io/repos/evalphobia/face-detect-annotator/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/face-detect-annotator?branch=master
[11]: https://codecov.io/github/evalphobia/face-detect-annotator/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/face-detect-annotator?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/face-detect-annotator
[14]: https://goreportcard.com/report/github.com/evalphobia/face-detect-annotator
[15]: https://img.shields.io/github/downloads/evalphobia/face-detect-annotator/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/face-detect-annotator/releases
[17]: https://img.shields.io/github/stars/evalphobia/face-detect-annotator.svg
[18]: https://github.com/evalphobia/face-detect-annotator/stargazers
[19]: https://codeclimate.com/github/evalphobia/face-detect-annotator/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/face-detect-annotator
[21]: https://bettercodehub.com/edge/badge/evalphobia/face-detect-annotator?branch=master
[22]: https://bettercodehub.com/

`face-detect-annotator` is a tool to detect faces from Image file by AWS Rekognition, Google Vision, Azure Computer Vision, OpenCV, Dlib, TensorFlow


# What's for?

Amazon advertise like "Our product is awesome!".

Google advocates like "Our AI is brilliant!!".

ML engineer says like "My top-noch model reach beyond the existing world!!!".


You may think, "Who is the winner?".
Okay, let's decide the No.1.


# Quick Usage

At first, install golang.
And install dependensies.

- OpenCV: https://gocv.io/getting-started/
- Dlib: https://github.com/Kagami/go-face
- TensorFlow: https://www.tensorflow.org/install/lang_go


Then, create binary file.

```bash
# Basic engines [pigo,azure,google,rekognition] (you can use these engines without OS dependencies or libraries)
$ go build -o face-detect-annotator ./cmd/face-detect-annotator-basic

# All engines
# $ go build -o face-detect-annotator ./cmd/face-detect-annotator-all
```

```bash
$ ./face-detect-annotator -h

Commands:

  help       show help
  detect     Detect faces from image file or csv list
  annotate   Annotate faces of image from --input TSV file
```



## Subcommands

### detect

`detect` command detecting faces from `--input` image file or csv list file.


```bash
$ ./face-detect-annotator detect -h

Detect faces from image file or csv list

Options:

  -h, --help                              display help information
  -i, --input                            *image dir path --input='/path/to/image_dir'
  -o, --output[=./output.tsv]            *output TSV file path --output='./output.tsv'
  -a, --all                               use all engines
  -e, --engine[=opencv,dlib,tensorflow]   comma separate Face Detect Engines --engine='opencv,dlib,tensorflow,rekognition,google,azure'
```

For example, if you want to detect faces of images from the CSV file,

```bash
cat ./input,csv

path,count
myimages/foobar/001.jpg,1
myimages/foobar/002.jpg,1
myimages/foobar/003.jpg,0
```

With all of the engines,

```bash
$ FDA_AZURE_REGION=westus \
FDA_AZURE_SUBSCRIPTION_KEY=ZZZ \
GOOGLE_APPLICATION_CREDENTIALS=$HOME/google/service_account.json \
AWS_ACCESS_KEY_ID=XXX \
AWS_SECRET_ACCESS_KEY=ZZZ \
./face-detect-annotator detector -i ./input.csv -o ./output.tsv  -e="rekognition,google,azure,tensorflow,opencv,dlib"

[INFO] Use rekognition
[INFO] Use google
[INFO] Use azure
[INFO] Use tensorflow
[INFO] Use opencv
[INFO] Use dlib
exec #: [2]
exec #: [1]
exec #: [0]
```

After a while, `output.tsv` will be created.


### annotate

`annotate` command draw rectangle lines around faces from TSV file, generated from `detect` command.


```bash
$ ./face-detect-annotator annotate -h

Annotate faces of image from --input TSV file

Options:

  -h, --help    display help information
  -i, --input  *detector's output tsv file --input='/path/to/output.tsv'
```

```bash
$ ./face-detect-annotator annotate -i ./output.tsv
```

After a while, `_annotated_` prefixed files will be created in the same directory of the given images.

```bash
$ tree

└── myimages
    └── foobar
        ├── _annotated_001.jpg
        ├── _annotated_002.jpg
        ├── _annotated_003.jpg
        ├── 001.jpg
        ├── 002.jpg
        └── 003.jpg
```

### Example output

![_annotated_01](https://user-images.githubusercontent.com/2827521/60887197-b3604f80-a28e-11e9-806c-c3c49f99211a.jpg)

![_annotated_02](https://user-images.githubusercontent.com/2827521/60887212-c115d500-a28e-11e9-8ce9-93063035cd23.jpg)


## Environment variables

| Name | Command | Description |
|:--|:--|:--|
| `FDA_ENGINE_ALL` | `detect` | Use all of the face detection engines. |
| `FDA_ENGINE_DLIB` | `detect` | Use Dlib. |
| `FDA_ENGINE_OPENCV` | `detect` | Use OpenCV. |
| `FDA_ENGINE_PIGO` | `detect` | Use Pigo. |
| `FDA_ENGINE_TF` | `detect` | Use TensorFlow. |
| `FDA_ENGINE_REKOGNITION` | `detect` | Use AWS Rekognition. |
| `FDA_ENGINE_GOOGLE` | `detect` | Use Google Vision API. |
| `FDA_ENGINE_AZURE` | `detect` | Use Azure Computer Vision API. |
| `FDA_DLIB_MODEL_DIR` | `detect` for Dlib | Specify the directory path of model files for Dlib. |
| `FDA_PIGO_CASCADE_FILE` | `detect` for Pigo | Specify the file path of a cascade file of Pigo. |
| `FDA_OPENCV_CASCADE_FILE` | `detect` for OpenCV | Specify the file path of a cascade file of OpenCV. |
| `FDA_TF_MODEL_FILE` | `detect` for TensorFloe | Specify the .pb file path of a model file for TensorFlow. |
| `FDA_AZURE_REGION` | `detect` for Azure | Specify the region for Azure. |
| `FDA_AZURE_SUBSCRIPTION_KEY` | `detect` for Azure | Specify the subscription key for Azure. |


# Credit

This library depends on these awesome libraries,

- Dlib: [github.com/Kagami/go-face](https://github.com/Kagami/go-face) by [Kagami](https://github.com/Kagami)
- OpenCV: [gocv.io/x/gocv](https://gocv.io/)
- Pigo: [github.com/esimov/pigo](https://github.com/esimov/pigo)
- TensorFlow: [github.com/tensorflow/tensorflow](https://github.com/tensorflow/tensorflow)
- Azure Computer Vision API: [github.com/Azure/azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go)
- Google Vision API: [github.com/googleapis/google-api-go-client](https://github.com/googleapis/google-api-go-client)
- AWS Rekognition: [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go)
