package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os/exec"
	"spirit2/websocket"
	"time"
)

const snapShotOrder = "screencap" //android2.3以上系统命令

func startTakeSnapshotJob() {
	s := websocket.GetNetStatus()
	if s {
		log.Printf("开始截图,地址[%s]", mac)
		t := time.Now()
		s := takeSnapshot()
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, s, &jpeg.Options{Quality: 10})
		data := map[string]string{
			"mac":     mac,
			"imgData": base64.StdEncoding.EncodeToString(buf.Bytes()),
		}
		bs, err := json.Marshal(data)
		url := "http://" + ip + ":" + port + "/updateScreenShot"
		request, err := http.NewRequest("POST", url, bytes.NewReader(bs))
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(t)
		log.Println("截图完成,耗时:", elapsed)
		var resp *http.Response
		resp, err = http.DefaultClient.Do(request)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
	} else {
		log.Println("网络连接失败,暂停截图")
	}
	time.AfterFunc(3*time.Second, startTakeSnapshotJob)
}

func takeSnapshot() (img *image.RGBA) {

	snapShotBuffer := new(bytes.Buffer)
	cmd := exec.Command(snapShotOrder)
	cmd.Stdout = snapShotBuffer
	if err := cmd.Run(); err != nil {
		return
	}
	var width, height, format int32
	binary.Read(snapShotBuffer, binary.LittleEndian, &width)
	binary.Read(snapShotBuffer, binary.LittleEndian, &height)
	err := binary.Read(snapShotBuffer, binary.LittleEndian, &format)
	if err != nil {
		return
	}
	w, h := int(width), int(height)
	//fmt.Println("宽:", w, "高:", h, "格式:", format)
	img = &image.RGBA{Pix: snapShotBuffer.Bytes(),
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	return
}
