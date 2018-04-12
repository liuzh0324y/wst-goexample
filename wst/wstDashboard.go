package wst

import (
	"sync"
	"time"
)

const maxErrLogLen = 128

type errEvent struct {
	Time time.Time `json:"t"`
	Err  string    `json:"e"`
}

type wstDashboard struct {
	lock sync.Mutex

	startTime time.Time

	totalWs      int
	totalRecvMsg int
	totalSendMsg int
	wsErrs       int
	httpErrs     int
}

type statusReport struct {
	UpTimeSec float64 `json:"upsec"`
	OpenWs    int     `json:"openws"`
	TotalWs   int     `json:"totalws"`
	WsErrs    int     `json:"wserrors"`
	HttpErrs  int     `json:"httperrors"`
}

func newDashboard() *wstDashboard {
	return &wstDashboard{startTime: time.Now()}
}

func (db *wstDashboard) getReport(rs *wstRoomTable) statusReport {
	db.lock.Lock()
	defer db.lock.Unlock()

	upTime := time.Since(db.startTime)
	return statusReport{
		UpTimeSec: upTime.Seconds(),
		OpenWs:    rs.wsCount(),
		TotalWs:   db.totalWs,
		WsErrs:    db.wsErrs,
		HttpErrs:  db.httpErrs,
	}
}

func (db *wstDashboard) incrWs() {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.totalWs += 1
}

func (db *wstDashboard) onHttpErr(err error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.httpErrs += 1
}
