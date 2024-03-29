package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/c1982/cloudping"
)

var (
	flagURL        = flag.String("url", "", "for custom test URL. Default nil")
	flagMethod     = flag.String("method", "GET", "HTTP method to use. Default GET")
	flagAWS        = flag.Bool("aws", false, "Test all AWS regions enpoints. Default false")
	flagShowBest   = flag.Bool("best", false, "Show only best region. Default false")
	flagCSV        = flag.Bool("csv", false, "Output CSV format. Default false")
	responseTmp    = `Region: %s	Total: %7dms	DNS Lookup: %7dms	TCP Connection: %7dms	First Byte Response: %7dms	Content Transfer: %7s`
	responseTmpCSV = `%s,%s,%s,%s,%s,%s`
	responseTmpErr = `Error: %v`
	printusage     = true
)

func main() {

	flag.Parse()

	cp := cloudping.CloudPing{}

	if *flagURL != "" {
		urlPing := cp.Ping(*flagMethod, *flagURL)
		write(urlPing)
		printusage = false
	}

	if *flagAWS {
		list := cp.RunAWSTestAsync()

		for i := 0; i < len(list); i++ {
			write(list[i])
			if *flagShowBest {
				if i == 0 {
					break
				}
			}
		}

		printusage = false
	}

	if printusage {
		usage()
	}
}

func write(p cloudping.PingItem) {

	if p.Err != nil {
		fmt.Printf(responseTmpErr, p.Err)
		return
	}

	if *flagCSV {
		fmt.Printf(responseTmpCSV+"\n",
			p.Region,
			p.TotalTime,
			p.DNSLookup,
			p.TCPConnection,
			p.FirstByteResponse,
			p.ContentTransfer)

	} else {

		fmt.Printf(responseTmp+"\n",
			p.Region,
			int(p.TotalTime/time.Millisecond),
			int(p.DNSLookup/time.Millisecond),
			int(p.TCPConnection/time.Millisecond),
			int(p.FirstByteResponse/time.Millisecond),
			p.ContentTransfer)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "EXAMPLE:")
	fmt.Fprintln(os.Stderr, "  cping -aws")
	fmt.Fprintln(os.Stderr, "  cping -aws -best -csv")
	fmt.Fprintln(os.Stderr, "  cping -url http://maestropanel.com")
}
