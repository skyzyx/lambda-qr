# QR Code Generator for AWS Lambda

Will take the request body that is POSTed to the endpoint, and generate a QR code from it.

This is a _Serverless_ app, written in Go ([Golang]), running in AWS Lambda, with API Gateway in front of it, and AWS CloudFront in front of that (for caching).

**This is an experiment.** Uptime is not guaranteed, and there is no SLA. _But_ all-in-all, it should be reasonably reliable.

> **NOTE:** I've tried to build this in a way that is very cheap to run. But if you find yourself using this more than just occasionally, consider kicking me down a few bucks to offset the cost of running this service. <https://cash.me/$rparman>

## Usage

The `https://qr.ryanparman.com` hostname is a CloudFront caching distribution, in front of API Gateway, in front of a Lambda function.

There are two endpoints:

* `qr.png` — This will return a (bitmap) PNG file.
* `qr.svg` — This will return a (vector) SVG file.

Additionally, it accepts 2 query-string parameters.

| Parameter | Example | Description |
| --------- | ------- | ----------- |
| `body` | `Hello world!` | (Required) The contents of the QR code. **No attempt is made to ascertain meaning from the input.** As such, it would be wise of you to review the QR Code formatting documentation at [github:zxing/zxing](https://github.com/zxing/zxing/wiki/Barcode-Contents). Value should be URL-encoded. |
| `size` | `300` | (Optional) For the PNG format, this is the length (in pixels) of one size of the square QR code. Allowed range is `150`–`1000`. The default value is `300`. |

> **NOTE:** Different browsers support different maximum lengths for URLs, with a generally-accepted answer of 2083 characters end-to-end. Some browsers support longer URLs, but in practice, trying to pack that much data into a little QR code results in a barely-usable QR code. As such, it is recommended that you stick to less verbose data structures.

### Response Bodies

The bodies of the responses contain PNG-formatted binary data, or SVG-formatted XML data. Both could be written directly to disk as a file.

### Example Data

#### vCard

Many apps understand the vCard specification. The wiki for [mangstadt/ez-vcard](https://github.com/mangstadt/ez-vcard) provides documentation for [Version differences](https://github.com/mangstadt/ez-vcard/wiki/Version-differences) and [Property lists](https://github.com/mangstadt/ez-vcard/wiki/Property-List) between the various versions of the vCard specification.

```vcard
BEGIN:VCARD
VERSION:3.0
N:Parman;Ryan;;;
FN:Ryan Parman
TITLE:Software/DevOps/Security Engineer
EMAIL;TYPE=INTERNET;TYPE=HOME;TYPE=pref:ryan@ryanparman.com
URL;TYPE=Homepage:https://ryanparman.com
URL;TYPE=GitHub:https://github.com/skyzyx
URL;TYPE=Keybase:https://keybase.io/skyzyx
X-SOCIALPROFILE;TYPE=twitter:https://twitter.com/skyzyx
END:VCARD
```

[Click](https://qr.ryanparman.com/qr.png?size=300&body=BEGIN%3AVCARD%0AVERSION%3A3.0%0AN%3AParman%3BRyan%3B%3B%3B%0AFN%3ARyan%20Parman%0ATITLE%3ASoftware%2FDevOps%2FSecurity%20Engineer%0AEMAIL%3BTYPE%3DINTERNET%3BTYPE%3DHOME%3BTYPE%3Dpref%3Aryan%40ryanparman.com%0AURL%3BTYPE%3DHomepage%3Ahttps%3A%2F%2Fryanparman.com%0AURL%3BTYPE%3DGitHub%3Ahttps%3A%2F%2Fgithub.com%2Fskyzyx%0AURL%3BTYPE%3DKeybase%3Ahttps%3A%2F%2Fkeybase.io%2Fskyzyx%0AX-SOCIALPROFILE%3BTYPE%3Dtwitter%3Ahttps%3A%2F%2Ftwitter.com%2Fskyzyx%0AEND%3AVCARD%0A)

#### Google Play Store

You can construct URIs that (on Android devices) link directly into Google Play Store. For example to encode a link to an app whose package is `com.agilebits.onepassword`, use:

```plain
market://details?id=com.agilebits.onepassword
```

[Click](https://qr.ryanparman.com/qr.png?size=300&body=market%3A%2F%2Fdetails%3Fid%3Dcom.agilebits.onepassword)

#### YouTube

You can use the `youtube:` protocol, or a standard link to YouTube.

```plain
youtube://ifwc5xgI3QM
https://youtu.be/ifwc5xgI3QM
https://www.youtube.com/watch?v=ifwc5xgI3QM
```

* [Apple URL Scheme Reference: YouTube Links](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/YouTubeLinks/YouTubeLinks.html#//apple_ref/doc/uid/TP40007899-CH8-SW1)

[Click](https://qr.ryanparman.com/qr.png?size=300&body=youtube%3A%2F%2Fifwc5xgI3QM)

## Developing/Deploying

### Golang

Go (when spoken) or [Golang] (when written) is a strongly-typed language from Google that "blends the simplicity of Python with the performance of C". Static binaries can be compiled for all major platforms, and many minor ones.

It is recommended that you install Golang using your system's package manager. If you don't have one (or if the version is too old), you can [install Golang from its website](https://golang.org/doc/install). Reading the [Getting Started](https://golang.org/doc/) documentation is a valuable exercise.

```bash
brew update && brew install golang
```

### Glide

Golang dependencies are managed with [Glide]. You should install them before compiling this project.

```bash
curl https://glide.sh/get | sh
glide install
```

### GoMetaLinter

[GoMetaLinter] pulls together many popular linting tools, and can run them on a project.

```bash
gometalinter.v2 --install
```

### Serverless

[Serverless] is a platform that wraps AWS Lambda and AWS CloudFormation, simplifying the deployment of Lambda apps. Serverless is written in Node.js, so you need to install that first.

I recommend you install the [Node Version Manager][nvm], and use that to install the latest Node.js and npm. Once that's complete, install `serverless`.

```bash
npm i -g serverless
```

### Developing

This app is small, and is self-contained in `main.go`.

_By default_, it expects to be running in AWS Lambda, receiving HTTP requests coming from API Gateway.

If you are performing local development/testing, run `make build` to build for the local platform, then `bin/qr -mock` to run it. The local build reads from `mock.json` and treats it as an incoming request from API Gateway. You can change the query-string parameters to have the app respond to the documented query string parameters.

```bash
make build

# PNG
bin/qr -mock

# SVG
QR_SVG=true bin/qr -mock
```

Make sure that you run the linter to catch any issues.

```bash
make lint
```

You can test some QR code data by reviewing [github:zxing/zxing](https://github.com/zxing/zxing/wiki/Barcode-Contents).

### Deployment

`serverless` uses the same [local credentials](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html) that the AWS CLI tools and the AWS SDKs use. If you haven't configured those yet, do that first.

Run `make package` to build a binary for AWS Lambda. Then, `serverless deploy` to deploy the app to your environment.

```bash
make package
serverless deploy
```

  [Glide]: https://glide.sh
  [Golang]: https://golang.org
  [GoMetaLinter]: https://github.com/alecthomas/gometalinter
  [nvm]: https://github.com/creationix/nvm
  [Serverless]: https://serverless.com/framework/docs/getting-started/
