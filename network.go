package codeutils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

func CallURL(req *http.Request, timeoutSec int) (result []byte, statusCode int, err error) {

	timeout := time.Duration(time.Duration(timeoutSec) * time.Second)

	client := http.Client{
		Timeout: timeout,
	}
	var response *http.Response
	response, err = client.Do(req)

	if err == nil {
		statusCode = response.StatusCode

		result, err = ioutil.ReadAll(response.Body)

	}

	return
}

func SetHeaderAuthentication(username, password string, req *http.Request) {

	req.SetBasicAuth(username, password)

}

func CallURLAsGet(url string, timeoutSec int) (result []byte, statusCode int, err error) {

	req, err := PrepareURLCall(url, "GET", nil)
	if err == nil {
		result, statusCode, err = CallURL(req, timeoutSec)
	}
	return
}

func CallURLAsPost(url string, contents []byte, timeoutSec int) (result []byte, statusCode int, err error) {

	req, err := PrepareURLCall(url, "POST", contents)
	if err == nil {
		result, statusCode, err = CallURL(req, timeoutSec)
	}
	return
}
