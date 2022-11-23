package codeutils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type CallURLResult struct {
	Content    []byte
	StatusCode int
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

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

}

func SetURLHeaders(req *http.Request, headers map[string]string) {

	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	return
}

func CallURL(req *http.Request, timeoutSec int) (result CallURLResult) {

	timeout := time.Duration(time.Duration(timeoutSec) * time.Second)

	client := http.Client{
		Timeout: timeout,
	}
	var response *http.Response
	response, result.Err = client.Do(req)

	if result.Err == nil {
		result.StatusCode = response.StatusCode

		result.Content, result.Err = ioutil.ReadAll(response.Body)

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
