package cloudping

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"time"
)

var (
	awsendpoints []string
)

func init() {
	awsendpoints = []string{"https://dynamodb.us-east-1.amazonaws.com/",
		"https://dynamodb.us-east-2.amazonaws.com/",
		"https://dynamodb.us-west-1.amazonaws.com/",
		"https://dynamodb.us-west-2.amazonaws.com/",
		"https://dynamodb.ca-central-1.amazonaws.com/",
		"https://dynamodb.eu-west-1.amazonaws.com/",
		"https://dynamodb.eu-west-2.amazonaws.com/",
		"https://dynamodb.eu-central-1.amazonaws.com/",
		"https://dynamodb.eu-west-3.amazonaws.com/",
		"https://dynamodb.eu-north-1.amazonaws.com/",
		"https://dynamodb.ap-south-1.amazonaws.com/",
		"https://dynamodb.ap-northeast-3.amazonaws.com/",
		"https://dynamodb.ap-northeast-2.amazonaws.com/",
		"https://dynamodb.ap-southeast-1.amazonaws.com/",
		"https://dynamodb.ap-southeast-2.amazonaws.com/",
		"https://dynamodb.ap-northeast-1.amazonaws.com/",
		"https://dynamodb.sa-east-1.amazonaws.com/",
		"https://dynamodb.cn-north-1.amazonaws.com.cn/",
		"https://dynamodb.cn-northwest-1.amazonaws.com.cn/",
		"https://dynamodb.us-gov-east-1.amazonaws.com/",
		"https://dynamodb.us-gov-west-1.amazonaws.com/",
	}
}

type CloudPing struct {
}

type PingItem struct {
	URL                 string
	Err                 error
	DNStime             time.Duration
	TCPtime             time.Duration
	FirstByteResponse   time.Duration
	ContentDownloadTime time.Duration
	ConnectionTime      time.Duration
}

func NewCloudPing() *CloudPing {
	return &CloudPing{}
}

func (c *CloudPing) RunAWSTest() []PingItem {

	list := []PingItem{}

	for i := 0; i < len(awsendpoints); i++ {
		var pi = PingItem{}
		pi, pi.Err = c.do("GET", awsendpoints[i])

		list = append(list, pi)
	}

	return list
}

func (c *CloudPing) do(method, url string) (ping PingItem, err error) {

	req, err := c.req(method, url)

	if err != nil {
		return ping, err
	}

	ping = PingItem{}
	var dnsStartTime, dnsDoneTime time.Time
	var gotConnTime, gotFRB time.Time
	var readBody time.Time

	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { dnsStartTime = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { dnsDoneTime = time.Now() },
		GotConn:              func(_ httptrace.GotConnInfo) { gotConnTime = time.Now() },
		GotFirstResponseByte: func() { gotFRB = time.Now() },
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	client := &http.Client{
		Transport: c.transport(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		return ping, err
	}

	resp.Body.Close()
	readBody = time.Now()

	if dnsStartTime.IsZero() {
		dnsStartTime = dnsDoneTime
	}

	ping.DNStime = dnsDoneTime.Sub(dnsStartTime)
	ping.TCPtime = gotConnTime.Sub(dnsDoneTime)
	ping.FirstByteResponse = gotFRB.Sub(gotConnTime)
	ping.ContentDownloadTime = readBody.Sub(gotFRB)
	ping.ConnectionTime = readBody.Sub(dnsStartTime)

	return ping, err
}

func (c *CloudPing) req(method, url string) (req *http.Request, err error) {
	return http.NewRequest(method, url, nil)
}

func (c *CloudPing) transport() *http.Transport {

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return tr
}
