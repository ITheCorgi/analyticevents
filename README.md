# Analytic Events

Fetch batch events from newline devided json POST request and stores into Clickhouse

Request:
``` 
POST /analytics/event/streaming
Headers: 	"Transfer-Encoding": "chunked"
	        "Connection": "keep-alive"
	        "Content-Type": "application/json"
Body:  newline devided json
```

```
goos: windows
goarch: amd64
pkg: github.com/ITheCorgi/analyticevents/tests
cpu: AMD Ryzen 5 4500U with Radeon Graphics
BenchmarkStream

BenchmarkStream-6            867           1744800 ns/op
PASS
```
