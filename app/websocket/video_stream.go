package websocket

import (
	"gocv.io/x/gocv"
	"github.com/yaphper/WebCamera/app/utils"
	"log"
	"image"
	"bytes"
	"image/jpeg"
	"time"
	"sync"
	"strconv"
	"runtime"
	"encoding/json"
	"strings"
	"github.com/Go-zh/net/websocket"
)

const (
	VideoCompressQuality 	= 40
	VideoWidth 				= 320
	VideoHeight 			= 240
	VideoFPS				= 21
)

type User struct {
	uid 			string
	ws				*websocket.Conn
	authenticated 	bool
}

var (
	videoCapture 		*gocv.VideoCapture
	frameBytes 			*bytes.Buffer

	frameWriteChan  	= make(chan *bytes.Buffer)
	closeFrameWriteChan = make(chan int)

	wsConnList 			= make(map[*User]*websocket.Conn)
	videoFrame 			= gocv.NewMat()

	wsConnListLock 		sync.Mutex
)

func init() {
	var err error
	videoCapture, err = gocv.OpenVideoCapture("/dev/video0")
	utils.CheckError(err)
	videoCapture.Set(gocv.VideoCaptureFrameWidth, VideoWidth)
	videoCapture.Set(gocv.VideoCaptureFrameWidth, VideoHeight)
	videoCapture.Set(gocv.VideoCaptureFPS, VideoFPS)
	videoCapture.Set(gocv.VideoCaptureBufferSize, 1)

	go videoStreamLoop()
	go writeFrameLoop()
}

func videoStreamLoop() {
	var (
		err error
		tempImage image.Image

		frameSize = image.Point{X: VideoWidth, Y: VideoHeight}
		compressOptions = jpeg.Options{Quality: VideoCompressQuality}
	)
	for {
		ok := videoCapture.Read(&videoFrame)
		if !ok {
			// close write loop
			closeFrameWriteChan <- 1
		}
		utils.CheckOk(ok, "failed to read camera")

		// resize
		gocv.Resize(videoFrame, &videoFrame, frameSize, 0, 0, gocv.InterpolationDefault)

		tempImage, err = videoFrame.ToImage()
		if err != nil {
			log.Print("failed to convert mat to image")
			continue
		}

		frameBytes = new(bytes.Buffer)
		jpeg.Encode(frameBytes, tempImage, &compressOptions)

		frameWriteChan <- frameBytes
	}
}

func writeFrameLoop() {
	var data *bytes.Buffer
	for {
		select {
		case data = <-frameWriteChan:
			go func() {
				for user := range wsConnList  {
					if user.authenticated {
						go writeFrame(user, data.Bytes())
					}
				}
			}()
		case <-closeFrameWriteChan:
			log.Print("frame write loop closed")
			goto ERR
		default:
			runtime.Gosched()
		}
	}

ERR:

}

func writeFrame(user *User, data []byte) {
	var (
		err error
	)

	conn, exists := wsConnList[user]
	if !exists {
		return
	}

	err = websocket.Message.Send(conn, data)
	if err != nil {
		log.Printf("socket[%s] write failed", conn.Request().RemoteAddr)

		conn.Close()
		log.Printf("socket[%s] closed", conn.Request().RemoteAddr)

		wsConnListLock.Lock()
		delete(wsConnList, user)
		wsConnListLock.Unlock()
	}
}

func heartbeatLoop(ws * websocket.Conn, user *User, closeChan chan int) {
	for {
		select {
		case <-closeChan:
			goto ERR
		default:
			err := websocket.Message.Send(ws, "heartbeat")
			if err != nil {
				log.Printf("socket[%s] heartbeat faild", ws.Request().RemoteAddr)
				goto ERR
			}
			time.Sleep(time.Second)
		}
	}

ERR:
	ws.Close()

	wsConnListLock.Lock()
	delete(wsConnList, user)
	wsConnListLock.Unlock()
}

func handleReceive(dataStr string, user *User) {
	var (
		userName 	string
		passWord	string
	)
	data := make(map[string]interface{})
	stringReader := strings.NewReader(dataStr)
	decoder := json.NewDecoder(stringReader)
	err := decoder.Decode(&data)
	if err != nil {
		println(err.Error())
	}
	userName = data["user"].(string)
	passWord = data["password"].(string)
	if CheckAccount(userName, passWord) {
		user.authenticated = true
		response(user.ws, "auth", "ok")
		log.Printf("user[%s] authed", userName)
	}
}

func VideoStreamWS(ws *websocket.Conn) {
	user := new(User)
	user.uid = strconv.FormatInt(time.Now().UnixNano(), 10)
	user.authenticated = false
	user.ws = ws

	wsConnListLock.Lock()
	wsConnList[user] = ws
	wsConnListLock.Unlock()

	// heartbeat loop
	closeChan := make(chan int)
	go heartbeatLoop(ws, user, closeChan)

	log.Printf("client[%s] connected", ws.Request().RemoteAddr)

	for {
		var readStream string
		err := websocket.Message.Receive(ws, &readStream)
		if err != nil {
			closeChan <- 1 // stop heartbeat loop
			break
		}
		if len(readStream) != 0 {
			go handleReceive(readStream, user)
		}
		//time.Sleep(time.Millisecond)
	}
}
