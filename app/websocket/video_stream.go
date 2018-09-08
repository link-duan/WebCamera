package websocket

import (
	"github.com/Go-zh/net/websocket"
	"gocv.io/x/gocv"
	"github.com/yaphper/WebCamera/app/utils"
	"log"
	image2 "image"
)

var videoCapture *gocv.VideoCapture

func init() {
	var err error
	videoCapture, err = gocv.OpenVideoCapture("/dev/video0")
	utils.CheckError(err)
	videoCapture.Set(gocv.VideoCaptureFrameWidth, 200)
	videoCapture.Set(gocv.VideoCaptureBufferSize, 1)
	videoCapture.Set(gocv.VideoCaptureFPS, 30)
}

func VideoStreamWS(ws *websocket.Conn) {
	defer ws.Close()

	var err error
	var encodedImage []byte
	image := gocv.NewMat()

	for {

		if !videoCapture.Read(&image) {
			log.Print("camera read failed")
			break
		}

		gocv.Resize(image, &image, image2.Point{200, 150}, 0, 0, gocv.InterpolationDefault)
		
		//gocv.CvtColor(image, &image, gocv.ColorRGBToGray)

		//gocv.Threshold(image, &image, 75, 255, gocv.ThresholdBinary)

		encodedImage, err = gocv.IMEncode(".jpg", image)

		err = websocket.Message.Send(ws, encodedImage)
		if err != nil {
			log.Print("socket write failed")
			break
		}
	}

}
