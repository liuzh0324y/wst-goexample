package wst

import (
	"encoding/json"
	"io"
	"log"

	"github.com/wisdomdev/wisdom-business-server/websocket"
)

type wsClientMsg struct {
	Cmd      string `json:"cmd"`
	RoomID   string `json:"roomid"`
	ClientID string `json:"clientid`
	Msg      string `json:"msg"`
}

type wsServerMsg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

type WstSignalServer struct {
}

func NewSignalServer() *WstSignalServer {
	obj := &WstSignalServer{}
	return obj
}

func (wst *WstSignalServer) Run() {

	// http.Handle("/ws", websocket.Handler(wst.wsHandler))

	// http.ListenAndServe(":8089", nil)
}

func (wst *WstSignalServer) wsHandler(ws *websocket.Conn) {
	log.Println("ws handler.")
	io.Copy(ws, ws)
}

// sendServerMsg sends a wsServerMsg composed from |msg| to the connection.
func sendServerMsg(w io.Writer, msg string) error {
	m := wsServerMsg{
		Msg: msg,
	}
	return send(w, m)
}

// sendServerErr sends a wsServerMsg composed from |errMsg| to the connection.
func sendServerErr(w io.Writer, errMsg string) error {
	m := wsServerMsg{
		Error: errMsg,
	}
	return send(w, m)
}

// Send writes a generic object as JSON to the writer.
func send(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}
