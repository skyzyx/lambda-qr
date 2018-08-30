# QR Code Generator for AWS Lambda

Will take the request body that is POSTed to the endpoint, and generate a QR code from it.

This is a _Serverless_ app, written in Go ([Golang]), running in AWS Lambda, with API Gateway in front of it, and AWS CloudFront in front of that (for caching).

**This is an experiment.** Uptime is not guaranteed, and there is no SLA. _But_ all-in-all, it should be reasonably reliable.

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

#### Plain Text

The most boring and unformatted of plain text.

```plain
Hello world!
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=Hello%20world!)

#### Wi-Fi connection

Wi-Fi connections are supported in iOS 11 (2017) and newer, and nearly all versions of Android going back to at least 2010.

```plain
WIFI:T:WPA;S:ThisIsMySSID;P:ThisIsMyPassword;;
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=WIFI%3AT%3AWPA%3BS%3AThisIsMySSID%3BP%3AThisIsMyPassword%3B%3B)

#### URL

Any URL is valid, but most QR readers will only support `http:` and `https`.

```plain
https://ryanparman.com
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=https%3A%2F%2Fryanparman.com)

#### Telephone

Standard `tel:` links should work here. See [CSS-Tricks: The Current State of Telephone Links](https://css-tricks.com/the-current-state-of-telephone-links/) and [Apple URL Scheme Reference: Phone Links](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/PhoneLinks/PhoneLinks.html) for some examples. You should generally use the most complete version of a telephone number possible (i.e., country code + area code + number).

```plain
# U.S. Directory assistance
tel:+18005551212
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=tel:+18005551212)

#### SMS/MMS/FaceTime

Similar to telephone links. See [CSS-Tricks: iPhone Calling and Texting Links](https://css-tricks.com/snippets/html/iphone-calling-and-texting-links/), [Apple URL Scheme Reference: SMS Links](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/SMSLinks/SMSLinks.html#//apple_ref/doc/uid/TP40007899-CH7-SW1), and [Apple URL Scheme Reference: FaceTime Links](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/FacetimeLinks/FacetimeLinks.html#//apple_ref/doc/uid/TP40007899-CH2-SW1) for some examples.

```plain
# Send an SMS/MMS to a number
sms:+18005551212

# Send an SMS/MMS to a number with pre-filled message.
sms:+18005551212:This%20is%20my%20text%20message.

# FaceTime Video
facetime:+18005551212
facetime:me@icloud.com

# FaceTime Audio
facetime-audio:+18005551212
facetime-audio:me@icloud.com
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=sms:+18005551212:This%20is%20my%20text%20message.)

#### Maps, Geo Coordinates

Geographic coordinates are as simple as the latitude + longitude.

```plain
geo:47.603363,-122.330417
```

[Click](https://qr.ryanparman.com/qr.png?size=150&body=geo:47.603363,-122.330417)

Services like [Apple Maps](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/MapLinks/MapLinks.html#//apple_ref/doc/uid/TP40007899-CH5-SW1) and [Google Maps](https://developers.google.com/maps/documentation/maps-static/intro) have more thorough implemenations with more options.

[Apple Maps](https://qr.ryanparman.com/qr.png?size=300&body=https%3A%2F%2Fmaps.apple.com%2F%3Faddress%3D400%2520Broad%2520St%2C%2520Seattle%2C%2520WA%2520%252098109%2C%2520United%2520States%26auid%3D17457489312301189071%26ll%3D47.620521%2C-122.349293%26lsp%3D9902%26q%3DSpace%2520Needle) [Google Maps](https://qr.ryanparman.com/qr.png?size=300&body=https%3A%2F%2Fmaps.google.com%2F%3Faddress%3D400%2520Broad%2520St%2C%2520Seattle%2C%2520WA%2520%252098109%2C%2520United%2520States%26auid%3D17457489312301189071%26ll%3D47.620521%2C-122.349293%26lsp%3D9902%26q%3DSpace%2520Needle)

#### Email

All of the standard `mailto:` tricks/links should work here as well. See [CSS-Tricks: Mailto Links](https://css-tricks.com/snippets/html/mailto-links/), [Apple URL Scheme Reference: Mail Links](https://developer.apple.com/library/archive/featuredarticles/iPhoneURLScheme_Reference/MailLinks/MailLinks.html#//apple_ref/doc/uid/TP40007899-CH4-SW1), and [RFC 6068](https://tools.ietf.org/html/rfc6068) for some examples.

```plain
# Address
mailto:someone@yoursite.com

# Address, subject
mailto:someone@yoursite.com?subject=Mail%20from%20Our%20Site

# Address, CC, BCC, subject
mailto:someone@yoursite.com?cc=someoneelse@theirsite.com,another@thatsite.com,me@mysite.com&bcc=lastperson@theirsite.com&subject=Big%20News

# Address, CC, BCC, subject, body
mailto:someone@yoursite.com?cc=someoneelse@theirsite.com,another@thatsite.com,me@mysite.com&bcc=lastperson@theirsite.com&subject=Big%20News&body=Body%20goes%20here.
```

[Click](https://qr.ryanparman.com/qr.png?size=300&body=mailto%3Asomeone%40yoursite.com%3Fcc%3Dsomeoneelse%40theirsite.com%2Canother%40thatsite.com%2Cme%40mysite.com%26bcc%3Dlastperson%40theirsite.com%26subject%3DBig%2520News%26body%3DBody%2520goes%2520here.)

#### Calendar Event

The newer [iCalendar](https://en.wikipedia.org/wiki/ICalendar) (`.ics`) format, as well as the older [vCalendar](https://en.wikipedia.org/wiki/ICalendar#vCalendar_1.0) (`.vcs`) format both define a sub-set of their formats for _events_. These are [vEvents](https://icalendar.org/iCalendar-RFC-5545/3-6-1-event-component.html).

```vevent
BEGIN:VEVENT
SUMMARY:Summer+Vacation!
DTSTART:20180601T070000Z
DTEND:20180831T070000Z
END:VEVENT
```

[Click](https://qr.ryanparman.com/qr.png?size=300&body=BEGIN%3AVEVENT%0ASUMMARY%3ASummer%2BVacation%21%0ADTSTART%3A20180601T070000Z%0ADTEND%3A20180831T070000Z%0AEND%3AVEVENT)

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

#### MECARD

[MECARD](https://en.wikipedia.org/wiki/MeCard_(QR_code)) is a format for contact information, which is much simpler (and less verbose) than vCard. It was created by NTT Docomo in Japan.

```mecard
MECARD:N:Parman,Ryan;EMAIL:ryan@ryanparman.com;URL:https://ryanparman.com;NOTE:https://github.com/skyzyx\nhttps://keybase.io/skyzyx\nhttps://twitter.com/skyzyx
```

[Click](https://qr.ryanparman.com/qr.png?size=300&body=MECARD%3AN%3AParman%2CRyan%3BEMAIL%3Aryan%40ryanparman.com%3BURL%3Ahttps%3A%2F%2Fryanparman.com%3BNOTE%3Ahttps%3A%2F%2Fgithub.com%2Fskyzyx%5Cnhttps%3A%2F%2Fkeybase.io%2Fskyzyx%5Cnhttps%3A%2F%2Ftwitter.com%2Fskyzyx)

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
