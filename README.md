# CloudPing [![Build Status](https://travis-ci.org/c1982/cloudping.svg?branch=master)](https://travis-ci.org/c1982/cloudping)

Cloud Ping is a command-line tool designed to estimate latency in AWS zones from your domain or server to which you are connected to the Internet.

[![Go Report Card](https://goreportcard.com/badge/github.com/c1982/cloudping)](https://goreportcard.com/report/github.com/c1982/cloudping)

### Download

[https://github.com/c1982/cloudping/releases](https://github.com/c1982/cloudping/releases)

## Usage

Check AWS endpoints latency of your location

```cping -aws```

Output:

```
Region: eu-central-1    Total:     584ms        DNS Lookup:     113ms   TCP Connection:     400ms       First Byte Response:      69ms  Content Transfer: 85.375µs
Region: eu-west-3       Total:     605ms        DNS Lookup:     127ms   TCP Connection:     396ms       First Byte Response:      81ms  Content Transfer: 64.27µs
Region: eu-north-1      Total:     612ms        DNS Lookup:     126ms   TCP Connection:     401ms       First Byte Response:      83ms  Content Transfer: 122.155µs
Region: eu-west-2       Total:     614ms        DNS Lookup:     127ms   TCP Connection:     402ms       First Byte Response:      84ms  Content Transfer: 73.99µs
Region: eu-west-1       Total:     648ms        DNS Lookup:     114ms   TCP Connection:     427ms       First Byte Response:     106ms  Content Transfer: 114.472µs
Region: us-east-1       Total:     759ms        DNS Lookup:     112ms   TCP Connection:     487ms       First Byte Response:     159ms  Content Transfer: 121.658µs
Region: ca-central-1    Total:     821ms        DNS Lookup:     147ms   TCP Connection:     504ms       First Byte Response:     169ms  Content Transfer: 73.301µs
Region: us-east-2       Total:     857ms        DNS Lookup:     138ms   TCP Connection:     538ms       First Byte Response:     180ms  Content Transfer: 91.123µs
Region: ap-south-1      Total:     880ms        DNS Lookup:     134ms   TCP Connection:     558ms       First Byte Response:     187ms  Content Transfer: 96.115µs
Region: us-gov-east-1   Total:     963ms        DNS Lookup:     214ms   TCP Connection:     570ms       First Byte Response:     178ms  Content Transfer: 102.894µs
Region: us-west-1       Total:    1057ms        DNS Lookup:     127ms   TCP Connection:     687ms       First Byte Response:     242ms  Content Transfer: 96.973µs
Region: us-west-2       Total:    1082ms        DNS Lookup:     114ms   TCP Connection:     721ms       First Byte Response:     245ms  Content Transfer: 94.287µs
Region: us-gov-west-1   Total:    1172ms        DNS Lookup:     136ms   TCP Connection:     790ms       First Byte Response:     246ms  Content Transfer: 108.831µs
Region: sa-east-1       Total:    1296ms        DNS Lookup:     113ms   TCP Connection:     889ms       First Byte Response:     293ms  Content Transfer: 113.744µs
Region: ap-southeast-1  Total:    1351ms        DNS Lookup:     136ms   TCP Connection:     914ms       First Byte Response:     300ms  Content Transfer: 48.492µs
Region: ap-northeast-3  Total:    1509ms        DNS Lookup:     136ms   TCP Connection:    1034ms       First Byte Response:     338ms  Content Transfer: 67.812µs
Region: cn-north-1      Total:    1521ms        DNS Lookup:     270ms   TCP Connection:     939ms       First Byte Response:     310ms  Content Transfer: 48.585µs
Region: ap-northeast-1  Total:    1535ms        DNS Lookup:     136ms   TCP Connection:    1051ms       First Byte Response:     347ms  Content Transfer: 48.814µs
Region: ap-northeast-2  Total:    1576ms        DNS Lookup:     126ms   TCP Connection:    1088ms       First Byte Response:     360ms  Content Transfer: 87.373µs
Region: ap-southeast-2  Total:    1730ms        DNS Lookup:     136ms   TCP Connection:    1186ms       First Byte Response:     407ms  Content Transfer: 76.36µs
Region: cn-northwest-1  Total:    1804ms        DNS Lookup:     471ms   TCP Connection:    1000ms       First Byte Response:     332ms  Content Transfer: 52.512µs
```

Show AWS best latency region

```cping -aws -best```

Output:

```
Region: eu-central-1    Total:     566ms        DNS Lookup:      95ms   TCP Connection:     392ms       First Byte Response:      77ms  Content Transfer: 88.74µs
```

Print CSV format

```cping -aws -csv```

Output:

```
eu-central-1,557.428034ms,86.129062ms,397.530328ms,73.61463ms,154.014µs
```

Test custom URL

```cping -url http://www.maestropanel.com```

```
Region: unknown Total:     167ms        DNS Lookup:      98ms   TCP Connection:      30ms       First Byte Response:      38ms  Content Transfer: 265.611µs
```

## Stats

| Stats | Description |
| ------ | ----------- |
| Total   | total time of all process |
| DNS Lookup | resolution time of the requested domain |
| TCP Connection | Time for TCP connection |
| First Byte Response | First byte of response from server |
| Content Transfer | Time to download the answer from the server |

## License

Distributed under the MIT License. See `LICENSE` for more information.


## Inspired

* https://cloudharmony.com/speedtest-for-aws
* http://www.cloudping.info/
* http://www.cloudwatch.in/
* https://cloudharmony.com/speedtest-for-aws

## Contact

Oğuzhan - [@c1982](https://twitter.com/c1982) - aspsrc@gmail.com

Project Link: [https://github.com/c1982/cloudping](https://github.com/c1982/cloudping)
