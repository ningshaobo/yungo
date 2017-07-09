package models

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	beegologs "github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 定义 wrap 结构体
type myTransport struct {
	// Uncomment this if you want to capture the transport
	// CapturedTransport http.RoundTripper
	//    http.RoundTripper

	Result interface{}
	body   []byte
}

/**
* @param  proxy 中间 wrap 处理
 */
func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	beegologs.Debug("路径（%v）进入 wrap", request.URL.Path)
	response, err := http.DefaultTransport.RoundTrip(request)
	// or, if you captured the transport
	// response, err := t.CapturedTransport.RoundTrip(request)

	// The httputil package provides a DumpResponse() func that will copy the
	// contents of the body into a []byte and return it. It also wraps it in an
	// ioutil.NopCloser and sets up the response to be passed on to the client.

	// You may want to check the Content-Type header to decide how to deal with
	// the body. In this case, we're assuming it's text.
	t.Result = response.StatusCode
	if t.Result.(int) < 400 {
		body, err := httputil.DumpResponse(response, true)
		if err != nil {
			// copying the response body did not work
			return nil, err
		}
		beegologs.Debug("clen= %v, blen = %v \n", response.ContentLength, len(body))
		bodyLen := int64(len(body)) - response.ContentLength
		t.body = body[bodyLen:]
	}

	return response, err
}

/**
* @param  反向代理 业务逻辑
 */
func ProxyHandler(inthis interface{}, myfun func(*http.Response) error, isResult bool) ([]byte, error) {
	remote, err := url.Parse("http://" + UtilsMasterProxyUrl())
	if err != nil {
		return nil, err
	}

	var this *beego.Controller
	this = inthis.(*beego.Controller)

	proxy := httputil.NewSingleHostReverseProxy(remote)

	if myfun != nil {
		beegologs.Debug("允许 modify body ")
		proxy.ModifyResponse = myfun
	}
	if isResult == true {
		var tt myTransport
		proxy.Transport = &tt // 是否加入proxy wrap 处理
		proxy.ServeHTTP(this.Ctx.ResponseWriter.ResponseWriter, this.Ctx.Request)
		beegologs.Debug("rsp status %v", tt.Result)
		if tt.Result.(int) < 400 {
			return tt.body, nil
		} else {
			err = fmt.Errorf("登录返回码小于 400，%v", tt.Result)
			return nil, err
		}
	} else {
		proxy.ServeHTTP(this.Ctx.ResponseWriter.ResponseWriter, this.Ctx.Request)
		return nil, nil
	}
}

/**
* @param  获取 response.Body
 */
func ProxyDrainBody(resp *http.Response) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error
	if resp.Body == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return &buf, nil
	}
	//	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return &buf, err
	}
	if err = resp.Body.Close(); err != nil {
		return &buf, err
	}
	return &buf, nil
}

/**
* @param  重新赋值 response body
 */
func ProxySetBody(resp *http.Response, buf *bytes.Buffer) error {
	// 重新定义 response 赋值
	resp.Body = ioutil.NopCloser(buf)
	resp.ContentLength = int64(buf.Len())
	resp.Header.Set("Content-Length", fmt.Sprint(buf.Len()))
	return nil
}
