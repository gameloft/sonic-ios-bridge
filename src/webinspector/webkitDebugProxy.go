package webinspector

import (
	"context"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var webDebug *WebkitDebugService
var localPort = 9222
var isAdapter = false

func SetIsAdapter(flag bool) {
	isAdapter = flag
}

func InitWebInspectorServer(udid string, port int, isProtocolDebug bool, isDTXDebug bool) context.CancelFunc {
	var err error
	var cannel context.CancelFunc
	if webDebug == nil {
		// 优化初始化过程
		ctx := context.Background()
		device := util.GetDeviceByUdId(udid)
		webDebug = NewWebkitDebugService(&device, ctx)
		cannel, err = webDebug.ConnectInspector()
		if err != nil {
			log.Fatal(err)
		}
	}
	localPort = port
	if isProtocolDebug {
		SetProtocolDebug(true)
	}
	if isDTXDebug {
		giDevice.SetDebug(true, true)
	}
	return cannel
}

func PagesHandle(c *gin.Context) {

	pages, err := webDebug.GetOpenPages(localPort)
	if err != nil {
		c.JSONP(http.StatusNotExtended, err)
	}
	c.JSONP(http.StatusOK, pages)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func PageDebugHandle(c *gin.Context) {
	id := c.Param("id")

	application, page, err := webDebug.FindPagesByID(id)
	if application == nil || page == nil {
		c.Error(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		log.Println(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		return
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	err = webDebug.StartCDP(application.ApplicationID, page.PageID, conn)
	if err != nil {
		log.Fatal(err)
	}

	//// 确保初始化完成
	if isAdapter {
		err = webDebug.ReceiveWebkitProtocolDataAdapter()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = webDebug.ReceiveWebkitProtocolData()
		if err != nil {
			fmt.Println(err)
		}
	}

	go func() {
		for {
			if isAdapter {
				err = webDebug.ReceiveWebkitProtocolDataAdapter()
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err = webDebug.ReceiveWebkitProtocolData()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
	if isAdapter {
		for {
			err = webDebug.ReceiveMessageToolAdapter()
			if err != nil {
				log.Panic(err)
			}
			if err == nil || err.Error() == "message is null" {
				continue
			}
		}
	} else {
		for {
			err = webDebug.ReceiveMessageTool()
			if err != nil {
				log.Panic(err)
			}
			if err == nil || err.Error() == "message is null" {
				continue
			}
		}
	}
}
