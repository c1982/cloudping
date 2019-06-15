package cloudping

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"sort"
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

type Pinglist []PingItem

func (pl Pinglist) Len() int           { return len(pl) }
func (pl Pinglist) Swap(i, j int)      { pl[i], pl[j] = pl[j], pl[i] }
func (pl Pinglist) Less(i, j int) bool { return pl[i].TotalTime < pl[j].TotalTime }

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

func (c *CloudPing) Ping(method, url string) PingItem {

	pi := c.do(method, url)
	pi.Region = "unknown"

	return pi
}

func (c *CloudPing) RunAWSTest() Pinglist {

	list := Pinglist{}

	for region, uri := range awsendpoints {

		pi := c.do("GET", uri)
		pi.Region = region

		list = append(list, pi)
	}

	sort.Sort(list)
	return list
}

func (c *CloudPing) RunAWSTestAsync() Pinglist {

	var wg sync.WaitGroup
	var items = Pinglist{}

	for rg, u := range awsendpoints {

		wg.Add(1)

		go func(region, uri string, wag *sync.WaitGroup) {

			pi := c.do("GET", uri)
			pi.Region = region

			items = append(items, pi)

			wag.Done()

		}(rg, u, &wg)
	}

	wg.Wait()

	sort.Sort(items)
	return items
}

func (c *CloudPing) do(method, url string) (ping PingItem) {

	req, err := c.req(method, url)

	if err != nil {
		ping.Err = err
		return ping
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
		ping.Err = err
		return ping
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

	return ping
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
