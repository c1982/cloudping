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

		fmt.Printf("ConnectionTime: %s\tTCPtime: %s\tDNS Time: %s\tFirst Byte: %s\tDownload: %s\n",
			p.ConnectionTime,
			p.TCPtime,
			p.DNStime,
			p.FirstByteResponse,
			p.ContentDownloadTime)
	}
}
