package apns

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	log4go "github.com/blackbeans/log4go"
	"golang.org/x/net/http2"
)

const (
	//开发环境
	URL_DEV = "api.development.push.apple.com:443"
	//正式环境
	URL_PRODUCTION = "api.push.apple.com:443"
)

type Notification struct {
	Topic       string
	ApnsID      string
	CollapseID  string
	Priority    int
	Expiration  time.Time
	DeviceToken string
	Payload     PayLoad
	Response    Response
}

//alert
type Alert struct {
	Body         string        `json:"body,omitempty"`
	ActionLocKey string        `json:"action-loc-key,omitempty"`
	LocKey       string        `json:"loc-key,omitempty"`
	LocArgs      []interface{} `json:"loc-args,omitempty"`
}

type Aps struct {
	Alert string `json:"alert,omitempty"`
	Badge int    `json:"badge,omitempty"` //显示气泡数
	Sound string `json:"sound"`           //控制push弹出的声音
}

type PayLoad struct {
	Aps Aps `json:"aps"`
}

//响应结果
type Response struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

type ApnsConn struct {
	ctx             context.Context
	cancel          context.CancelFunc
	cert            *tls.Config // ssl 证书
	hostport        string
	worktime        time.Time
	keepalivePeriod time.Duration

	c     *http2.ClientConn
	conn  net.Conn
	alive bool
}

//NewApnsConn ...
func NewApnsConn(
	ctx context.Context,
	certificates tls.Certificate,
	hostport string,
	keepalivePeriod time.Duration) (*ApnsConn, error) {

	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{certificates}
	tlsConfig.InsecureSkipVerify = true
	if len(certificates.Certificate) > 0 {
		tlsConfig.BuildNameToCertificate()
	}

	conn := &ApnsConn{
		ctx:             ctx,
		cert:            tlsConfig,
		hostport:        hostport,
		keepalivePeriod: keepalivePeriod}
	err := conn.Open()
	go conn.keepalive()
	return conn, err
}

func (apnsConn *ApnsConn) keepalive() {
	ticker := time.NewTicker(5 * time.Second)
	for apnsConn.alive {
		select {
		case <-ticker.C:
			if nil != apnsConn.c && apnsConn.alive && time.Since(apnsConn.worktime) > selc.keepalivePeriod {
				err := apnsConn.c.Ping(apnsConn.ctx)
				if err != nil {
					log4go.WarnLog("service", "CheckAlive|%s|Ping|FAIL|...", apnsConn.hostport)
					apnsConn.close0()
					//重新连接
					apnsConn.Open()
				} else {
					log4go.DebugLog("service", "CheckAlive|%s|Ping|SUCC|...", apnsConn.hostport)
					break
				}
			}
		case <-apnsConn.ctx.Done():
			ticker.Stop()
			if apnsConn.c != nil {
				apnsConn.Destory()
			}
		}
	}
}

func (apnsConn *ApnsConn) Open() error {
	dialer := &net.Dialer{
		Timeout:   apnsConn.keepalivePeriod * 2,
		KeepAlive: apnsConn.keepalivePeriod,
	}

	DialTLS := func(network, addr string, cfg *tls.Config) (net.Conn, error) {
		conn, err := tls.DialWithDialer(dialer, network, apnsConn.hostport, apnsConn.cert)
		if err != nil {
			return nil, err
		}
		return conn, err
	}
	conn, err := DialTLS("tcp", apnsConn.hostport, apnsConn.cert)
	if err != nil {
		return err
	}
	transport := &http2.Transport{
		TLSClientConfig: apnsConn.cert,
	}
	h2c, err := transport.NewClientConn(conn.(*tls.Conn))
	if err != nil {
		return err
	}

	apnsConn.c = h2c
	apnsConn.conn = conn
	apnsConn.alive = true
	log.Printf("Reconnect Apns|SUCC|...")
	return nil
}

func (apnsConn *ApnsConn) SendMessage(notification *Notification) error {
	data, err := json.Marshal(notification.Payload)
	if nil != err {
		return errors.New("Invalid Payload !")
	}

	domain, _, _ := net.SplitHostPort(apnsConn.hostport)
	url := fmt.Sprintf("https://%s/3/device/%v", domain, notification.DeviceToken)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("CreateReq|FAIL|%v|%s|%s", err, url, string(data))
		return err
	}
	setHeaders(req, notification)
	response, err := apnsConn.c.RoundTrip(req)
	if err != nil {
		log.Printf("FireReq|FAIL|%v|%s|%s", err, url, string(data))
		return err
	}
	defer response.Body.Close()

	apnsConn.worktime = time.Now()
	resp := &Response{}
	resp.Status = response.StatusCode
	if err = decoder.Decode(&resp); nil != err && err != io.EOF {
		log.Printf("UnmarshaldResponse|FAIL|%v|%s|%s", err, url, string(data))
		return err
	}
	notification.Response = *resp
	return nil

}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}
	if n.CollapseID != "" {
		r.Header.Set("apns-collapse-id", n.CollapseID)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", fmt.Sprintf("%v", n.Priority))
	}
	if !n.Expiration.IsZero() {
		r.Header.Set("apns-expiration", fmt.Sprintf("%v", n.Expiration.Unix()))
	}
}
