package main

import (
	"cloudping"
	"flag"
	"fmt"
	"os"
	"time"
)

var flagURL = flag.String("url", "", "for custom test URL. Default nil")
var flagMethod = flag.String("method", "GET", "HTTP method to use. Default GET")
var flagAWS = flag.Bool("aws", false, "Test all AWS regions enpoints. Default false")
var flagShowBest = flag.Bool("best", false, "Show only best region. Default false")
var flagCSV = flag.Bool("csv", false, "Output CSV format. Default false")

var responseTmp = `Region: %s	Total: %7dms	DNS Lookup: %7dms	TCP Connection: %7dms	First Byte Response: %7dms	Content Transfer: %7s`
var responseTmpCSV = `%s,%s,%s,%s,%s,%s`

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
		os.Exit(2)
	}

	cp := cloudping.CloudPing{}

	if *flagURL != "" {
		urlPing := cp.Ping(*flagMethod, *flagURL)
		write(urlPing)
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
	}
}

func write(p cloudping.PingItem) {

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
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "EXAMPLE:")
	fmt.Fprintln(os.Stderr, "  cp -aws")
	fmt.Fprintln(os.Stderr, "  cp -aws -best -csv")
	fmt.Fprintln(os.Stderr, "  cp -url http://maestropanel.com")
}
