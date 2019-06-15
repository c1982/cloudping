package cloudping

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

var (
	awsendpoints map[string]string
)

func init() {
	awsendpoints = map[string]string{
		"us-east-1":      "https://dynamodb.us-east-1.amazonaws.com/",
		"us-east-2":      "https://dynamodb.us-east-2.amazonaws.com/",
		"us-west-1":      "https://dynamodb.us-west-1.amazonaws.com/",
		"us-west-2":      "https://dynamodb.us-west-2.amazonaws.com/",
		"ca-central-1":   "https://dynamodb.ca-central-1.amazonaws.com/",
		"eu-west-1":      "https://dynamodb.eu-west-1.amazonaws.com/",
		"eu-west-2":      "https://dynamodb.eu-west-2.amazonaws.com/",
		"eu-central-1":   "https://dynamodb.eu-central-1.amazonaws.com/",
		"eu-west-3":      "https://dynamodb.eu-west-3.amazonaws.com/",
		"eu-north-1":     "https://dynamodb.eu-north-1.amazonaws.com/",
		"ap-south-1":     "https://dynamodb.ap-south-1.amazonaws.com/",
		"ap-northeast-3": "https://dynamodb.ap-northeast-3.amazonaws.com/",
		"ap-northeast-2": "https://dynamodb.ap-northeast-2.amazonaws.com/",
		"ap-southeast-1": "https://dynamodb.ap-southeast-1.amazonaws.com/",
		"ap-southeast-2": "https://dynamodb.ap-southeast-2.amazonaws.com/",
		"ap-northeast-1": "https://dynamodb.ap-northeast-1.amazonaws.com/",
		"sa-east-1":      "https://dynamodb.sa-east-1.amazonaws.com/",
		"cn-north-1":     "https://dynamodb.cn-north-1.amazonaws.com.cn/",
		"cn-northwest-1": "https://dynamodb.cn-northwest-1.amazonaws.com.cn/",
		"us-gov-east-1":  "https://dynamodb.us-gov-east-1.amazonaws.com/",
		"us-gov-west-1":  "https://dynamodb.us-gov-west-1.amazonaws.com/",
	}
}

type CloudPing struct {
}

type PingItem struct {
	Region            string
	URL               string
	Err               error
	DNSLookup         time.Duration
	TCPConnection     time.Duration
	FirstByteResponse time.Duration
	ContentTransfer   time.Duration
	TotalTime         time.Duration
}

func NewCloudPing() *CloudPing {
	return &CloudPing{}
}

func (c *CloudPing) RunAWSTest() []PingItem {

	list := []PingItem{}

	for region, uri := range awsendpoints {

		var pi = PingItem{}
		pi.Region = region
		pi, pi.Err = c.do("GET", uri)

		list = append(list, pi)
	}

	return list
}

func (c *CloudPing) RunAWSTestAsync() []PingItem {

	var wg sync.WaitGroup
	var items = []PingItem{}

	for rg, u := range awsendpoints {

		wg.Add(1)

		go func(region, uri string, wag *sync.WaitGroup) {

			pi, err := c.do("GET", uri)
			pi.Err = err
			pi.Region = region

			items = append(items, pi)

			wag.Done()

		}(rg, u, &wg)
	}

	wg.Wait()

	return items
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

	ping.DNSLookup = dnsDoneTime.Sub(dnsStartTime)
	ping.TCPConnection = gotConnTime.Sub(dnsDoneTime)
	ping.FirstByteResponse = gotFRB.Sub(gotConnTime)
	ping.ContentTransfer = readBody.Sub(gotFRB)
	ping.TotalTime = readBody.Sub(dnsStartTime)

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
