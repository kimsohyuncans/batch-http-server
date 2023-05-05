### Batch the expensive operations in API

#### prerequisite
[vegeta](https://github.com/tsenart/vegeta)
[golang](https://golang.org/)

### How To Run

run the http server
```bash
go run main.go
```

simulate the load test
```bash
echo "GET http://localhost:1323/summary" | vegeta attack -duration=30s -rate=5 -timeout=3s | vegeta report
```

you can monitor the proccess in "http://localhost:1323/stats"
