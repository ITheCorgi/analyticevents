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
