package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aaronarduino/goqrsvg"
	"github.com/ajstarks/svgo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	// "github.com/davecgh/go-spew/spew"
)

var (
	buf        bytes.Buffer
	emptyAGW   *events.APIGatewayProxyResponse
	err        error
	isMock     *bool
	isSVG      bool
	mimetype   string
	size       int64
	statusCode int
)

// The API Gateway handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	statusCode = int(200)
	emptyAGW = new(events.APIGatewayProxyResponse)

	if request.QueryStringParameters["body"] == "" {
		return *emptyAGW, errors.New("The 'body' query string parameter is required.")
	}

	// Create the barcode
	qrCode, _ := qr.Encode(request.QueryStringParameters["body"], qr.H, qr.Auto)

	// Do we want an SVG?
	if isSVG, err = strconv.ParseBool(os.Getenv("QR_SVG")); err != nil {
		isSVG = false
	}

	if isSVG {
		// SVG
		mimetype = "image/svg+xml"
		s := svg.New(&buf)

		// Write QR code to SVG
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(s)
		err = qs.WriteQrSVG(s)

		if err != nil {
			return *emptyAGW, err
		}

		s.End()
	} else {
		// PNG
		mimetype = "image/png"

		// What size should we create?
		if size, err = strconv.ParseInt(request.QueryStringParameters["size"], 10, 64); err != nil {
			size = 300
		}

		// Cap the size of the PNG image
		if size > 1000 {
			size = 1000
		}
		if size < 150 {
			size = 150
		}

		// Scale the barcode to size pixels
		qrCode, err = barcode.Scale(qrCode, int(size), int(size))

		if err != nil {
			return *emptyAGW, err
		}

		err = png.Encode(&buf, qrCode)

		if err != nil {
			return *emptyAGW, err
		}
	}

	cacheFrom := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(1, 0, 0).Format(http.TimeFormat)

	imageBinary := buf.Bytes()
	buf.Reset()

	// API Gateway has this weird requirement for Lambda where you can't just return data without it being corrupted.
	// Instead, you need to Base64-encode it coming out of Lambda, then tell API Gateway that it is Base64-encoded,
	// which will decode it on the pass back to the client.
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  mimetype,
			"Last-Modified": cacheFrom,
			"Expires":       cacheUntil,
		},
		Body:            base64.StdEncoding.EncodeToString(imageBinary),
		StatusCode:      statusCode,
		IsBase64Encoded: true,
	}, nil
}

// The core function
func main() {
	isMock = flag.Bool("mock", false, "Read from the local `mock.json` file instead of an API Gateway request.")
	flag.Parse()

	if *isMock {
		// read json from file
		inputJSON, jsonErr := ioutil.ReadFile("./mock.json")
		if jsonErr != nil {
			fmt.Println(jsonErr.Error())
			os.Exit(1)
		}

		// de-serialize into Go object
		var inputEvent events.APIGatewayProxyRequest
		if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		response, err := Handler(inputEvent)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		data, err := base64.StdEncoding.DecodeString(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(data))
	} else {
		lambda.Start(Handler)
	}
}
