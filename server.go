package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"spirit2/websocket"
	"strings"
	"time"
)

func main() {
	var (
		err  error
		data []byte
	)
	h := websocket.Handler(func(conn *websocket.Connection) {
		for {
			if data, err = conn.ReadMessage(); err != nil {
				goto ERR
			}
			log.Printf("收到 message:[%s]", string(data))
		}
	ERR:
		//TODO :关闭连接
		conn.Close()
	})
	go func() {
		fmt.Printf("启动连接检查器......\n")
		for {
			fmt.Printf("********************* 连接数目:[%d] *************************\n", len(websocket.WsMap.Map))
			for k, v := range websocket.WsMap.Map {
				fmt.Printf("serverId:[%s] <---------> clientId:[%s]\n", k, v)
			}
			if len(websocket.WsMap.Map) > 0 {
				fmt.Println("********************************************************************************************")
			}

			time.Sleep(5 * time.Second)
		}

	}()
	http.Handle("/ws", h)
	http.HandleFunc("/updateScreenShot", updateScreenShot)
	http.ListenAndServe(":8080", nil)
}

func updateScreenShot(response http.ResponseWriter, r *http.Request) {
	t := time.Now()
	log.Println("服务端收到数据.....")
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	json.Unmarshal(bs, &data)
	log.Printf("截图请求来自:物理地址[%s]", data["mac"].(string))

	deviceMark := strings.ToUpper(data["mac"].(string))
	imgDataStr := data["imgData"].(string)

	folderPath := "./snapshot/" + strings.Replace(deviceMark, ":", "--", -1)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	out, _ := os.Create(folderPath + "/t11.png")
	defer out.Close()
	imgDataBs, _ := base64.StdEncoding.DecodeString(imgDataStr)
	out.Write(imgDataBs)
	elapsed := time.Since(t)
	log.Println("服务端渲染完成,耗时:", elapsed)
	log.Println("----------------------------------------------------")

}
