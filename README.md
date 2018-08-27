# QR Code Generator for AWS Lambda

> **IMPORTANT:** Work-in-Progress. This package does not (yet) work the way it is described below.

Will take the request body that is POSTed to the endpoint, and generate a QR code from it.

This is a _Serverless_ app, written in Go ([Golang]), running in AWS Lambda, with API Gateway in front of it, and AWS CloudFront in front of that (for caching).

**This is an experiment.** Uptime is not guaranteed, and there is no SLA. _But_ all-in-all, it should be reasonably reliable.

## Usage

The `{endpoint}` endpoint is an API Gateway configuration sitting in front of a Lambda function.

The contents of the QR code should be sent as the request body in a `POST` request to the endpoint. **No attempt is made to ascertain meaning from the input.** As such, it would be wise of you to review the QR Code formatting documentation at [github:zxing/zxing](https://github.com/zxing/zxing/wiki/Barcode-Contents).

Additionally, it accepts 2 query-string parameters.

| Parameter | Example | Description |
| --------- | ------- | ----------- |
| `svg` | `true`\|`false` | (Optional) Whether or not an SVG response should be returned. A value of `true` means that the response will be in SVG format. A value of `false` means that the response will be in PNG format. The default value is `false`. |
| `size` | `300` | (Optional) For the PNG format, this is the length (in pixels) of one size of the square QR code. The default value is `300`. Any values larger than `1000` will be rounded down to `1000`. |

```
https://{hostname}/dev/qr
https://{hostname}/dev/qr?size=300
https://{hostname}/dev/qr?svg=true
```

### Response Bodies

The bodies of the responses contain PNG-formatted binary data, or SVG-formatted XML data. Both could be written directly to disk as a file.

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
bin/qr -mock
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
