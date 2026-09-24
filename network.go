package codeutils

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type CallURLResult struct {
	Content    []byte
	StatusCode int
	Status     string
	Err        error
}

func GetRemoteIP(r *http.Request) (ip string) {

	ip = r.RemoteAddr
	if !strings.Contains(ip, "::") {
		ip = ip[:strings.Index(ip, ":")]
	}
	if strings.HasPrefix(ip, "127.") || strings.HasPrefix(ip, "0.") ||
		strings.HasPrefix(ip, "0:") || strings.HasPrefix(r.RemoteAddr, "[::1]") {
		ipt := r.Header.Get("X-Forwarded-For")
		if ipt != "" {
			ip = ipt
		} else {
			ipt = r.Header.Get("X-Real-IP")
			if ipt != "" {
				ip = ipt
			}
		}
	}

	return
}

func PrepareURLCall(url string, method string, content []byte) (req *http.Request, err error) {

	method = strings.ToUpper(method)
	if method == "" {
		method = "GET"
	}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
	if err == nil {
		req.Close = true
	}

	return
}

func SkipHTTPSVerfication() {
	defaultTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// prefer : 250, 10 , 5
func ApplyTransportHttp2Options(maxConcurStream int, pingTimeout int, SendingPingTimeout int, ForceHTTP2 bool){
		defaultTransport.ForceAttemptHTTP2=   ForceHTTP2
		defaultTransport.HTTP2= &http.HTTP2Config{
		  MaxConcurrentStreams: maxConcurStream,
			SendPingTimeout: time.Second * time.Duration(SendingPingTimeout),
			PingTimeout:     time.Second * time.Duration(pingTimeout),
	}
	
}
// prefer : 100, 10, 20, 5, 5
func ApplyTransportConnectionsOptions(maxIdleConn, maxIdleConnPerHost, IdleConnTimeout, dialerTimeout, dialerkeepalive int )  {
	defaultTransport.MaxIdleConns = maxIdleConn
	defaultTransport.MaxIdleConnsPerHost = maxIdleConnPerHost
	defaultTransport.IdleConnTimeout = time.Duration(IdleConnTimeout) * time.Second

	dialer := &net.Dialer{
		Timeout:   time.Duration(dialerTimeout) * time.Second,
		KeepAlive: time.Second * time.Duration(dialerkeepalive),
	}
	 defaultTransport.DialContext = dialer.DialContext
	
}

func SetURLHeaders(req *http.Request, headers map[string]string) {

	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	return
}

var defaultTransport = http.DefaultTransport.(*http.Transport).Clone()
var defaultClient = &http.Client{
	Transport: defaultTransport,
}

func CallURL(req *http.Request, timeoutSec int) (result CallURLResult) {

	timeout := time.Duration(time.Duration(timeoutSec) * time.Second)
	ctx, cancel := context.WithTimeout(req.Context(), timeout)
	defer cancel()

	req = req.WithContext(ctx)
	var response *http.Response
	response, result.Err = defaultClient.Do(req)

	if result.Err == nil {
		defer response.Body.Close()
		result.StatusCode = response.StatusCode
		result.Status = response.Status

		result.Content, result.Err = io.ReadAll(response.Body)

	} else {
		result.StatusCode = http.StatusInternalServerError

	}

	return
}

func SetHeaderAuthentication(req *http.Request, username, password string) {

	req.SetBasicAuth(username, password)

}

func CallURLAsGet(url string, timeoutSec int) (result CallURLResult) {

	var req *http.Request
	req, result.Err = PrepareURLCall(url, "GET", nil)
	if result.Err == nil {
		result = CallURL(req, timeoutSec)
	}
	return
}

func CallURLAsPost(url string, contents []byte, timeoutSec int) (result CallURLResult) {

	var req *http.Request

	req, result.Err = PrepareURLCall(url, "POST", contents)
	if result.Err == nil {
		result = CallURL(req, timeoutSec)
	}
	return
}
