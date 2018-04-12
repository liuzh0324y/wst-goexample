package wst

import (
	"log"
	"net/http"
	"time"

	"github.com/wisdomdev/wisdom-business-server/websocket"
)

const registerTimeoutSec = 10
const wsReadTimeOutSec = 60 * 60 * 24

type WstHttpServer struct {
	addr     string
	path     string
	certFile string
	keyFile  string

	*wstRoomTable
	dash *wstDashboard
}

func NewHttpServer() *WstHttpServer {
	obj := &WstHttpServer{
		addr:     ":8090",
		path:     "/home/liuzh/work/workspace/src/github.com/wisdomdev/wisdom-business-server/static/html",
		certFile: "/home/liuzh/work/workspace/src/github.com/wisdomdev/wisdom-business-server/static/key/cert.pem",
		keyFile:  "/home/liuzh/work/workspace/src/github.com/wisdomdev/wisdom-business-server/static/key/key.pem",

		wstRoomTable: newRoomTable(time.Second*registerTimeoutSec, "wst"),
		dash:         newDashboard(),
	}
	return obj
}

func (wst *WstHttpServer) Run() {
	http.HandleFunc("/", wst.rootHandler)
	// http.Handle("/js/", http.FileServer(http.Dir(wst.path+"/js")))
	// http.Handle("/css/", http.FileServer(http.Dir(wst.path+"/css")))
	http.Handle("/wst", websocket.Handler(wst.wsHandler))

	err := http.ListenAndServeTLS(wst.addr, wst.certFile, wst.keyFile, nil)
	// err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (wst *WstHttpServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(wst.path)).ServeHTTP(w, r)
}

func (wst *WstHttpServer) wsHandler(ws *websocket.Conn) {
	var rid, cid string

	registered := false

	var msg wsClientMsg

loop:
	for {
		err := ws.SetReadDeadline(time.Now().Add(time.Duration(wsReadTimeOutSec) * time.Second))
		if err != nil {
			log.Fatal("ws.SetReadDeadline error: "+err.Error(), ws)
			break
		}

		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal("websocket.JSON.Receive error: "+err.Error(), ws)
			}
			break
		}

		switch msg.Cmd {
		case "register":
			if registered {
				log.Fatal("Duplicated register request", ws)
				break loop
			}
			if msg.RoomID == "" || msg.ClientID == "" {
				log.Fatal("Invalid register request: missing 'clientid' or 'roomid'", ws)
				break loop
			}
			if err = wst.wstRoomTable.register(msg.RoomID, msg.ClientID, ws); err != nil {
				log.Fatal(err.Error(), ws)
				break loop
			}
			registered, rid, cid = true, msg.RoomID, msg.ClientID
			wst.dash.incrWs()

			defer wst.wstRoomTable.deregister(rid, cid)
			break

		case "send":
			if !registered {
				log.Fatal("Client not registered", ws)
				break loop
			}
			if msg.Msg == "" {
				log.Fatal("Invalid send request: missing 'msg'", ws)
				break loop
			}
			wst.wstRoomTable.send(rid, cid, msg.Msg)
			break

		default:
			log.Fatal("Invalid message: unexpected 'cmd'", ws)
			break
		}
	}

	ws.Close()
}
