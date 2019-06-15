package cloudping

import (
	"fmt"
	"testing"
)

func Test_Ping(t *testing.T) {

	cp := NewCloudPing()
	list := cp.RunAWSTest()

	if len(list) == 0 {
		t.Error("list is empty!")
	}

	for i := 0; i < len(list); i++ {
		p := list[i]

		if p.Err != nil {
			t.Error(p.Err)
		}

		fmt.Printf("Region: %s\tTotal: %s\tTCPtime: %s\tDNS Time: %s\tFirst Byte: %s\tDownload: %s\n",
			p.Region,
			p.TotalTime,
			p.TCPConnection,
			p.DNSLookup,
			p.FirstByteResponse,
			p.ContentTransfer)
	}
}

func Test_Ping_Async(t *testing.T) {

	cp := NewCloudPing()
	list := cp.RunAWSTestAsync()

	for i := 0; i < len(list); i++ {

		p := list[i]

		if p.Err != nil {
			t.Error(p.Err)
		}

		fmt.Printf("Region: %s\tTotal: %s\tTCPtime: %s\tDNS Time: %s\tFirst Byte: %s\tDownload: %s\n",
			p.Region,
			p.TotalTime,
			p.TCPConnection,
			p.DNSLookup,
			p.FirstByteResponse,
			p.ContentTransfer)
	}
}
