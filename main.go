package main

import (
	"bytes"
	"encoding/json"
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
)

var (
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

	// Create the barcode
	qrCode, _ := qr.Encode(request.Body, qr.H, qr.Auto)

	// Do we want an SVG?
	if isSVG, err = strconv.ParseBool(request.QueryStringParameters["svg"]); err != nil {
		isSVG = false
	}

	// Write image data to buffer
	buf := new(bytes.Buffer)

	if isSVG {
		// SVG
		mimetype = "image/svg+xml"
		s := svg.New(buf)

		// Write QR code to SVG
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(s)
		err = qs.WriteQrSVG(s)

		if err != nil {
			panic(err)
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

		// Scale the barcode to size pixels
		qrCode, err = barcode.Scale(qrCode, int(size), int(size))

		if err != nil {
			panic(err)
		}

		err = png.Encode(buf, qrCode)

		if err != nil {
			panic(err)
		}
	}

	cacheFrom := time.Now().Format(http.TimeFormat)

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":  mimetype,
			"Last-Modified": cacheFrom,
			"Expires":       cacheFrom,
		},
		Body:       buf.String(),
		StatusCode: statusCode,
	}, nil
}

// The core function
func main() {
	isMock = flag.Bool("mock", false, "Read from the local `mock.json` file instead of an API Gateway request.")
	flag.Parse()

	if *isMock {
		// read json from file
		inputJSON, err := ioutil.ReadFile("./mock.json")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// de-serialize into Go object
		var inputEvent events.APIGatewayProxyRequest
		if err := json.Unmarshal(inputJSON, &inputEvent); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		response, _ := Handler(inputEvent)

		fmt.Println(response.Body)
	} else {
		lambda.Start(Handler)
	}
}
