# gosd

**gosd** is a Go Runtime stats exporter into Stackdriver.

**gosd** can monitor these stats to better understand the overall health and performance of Go programs.

## Monitored Stats

- [runtime.NumGoroutine](https://golang.org/pkg/runtime/debug/#Stack)
- [runtime.NumCgoCall](https://golang.org/pkg/runtime/#NumCgoCall)
- [runtime.ReadMemStats](https://golang.org/pkg/runtime/#ReadMemStats)

## Setting up authentication

To run the library, you must first set up authentication by creating a service account and setting an environment variable.

[Create service account](https://console.cloud.google.com/apis/credentials/serviceaccountkey) with write acccess for Monitoring, download the key and set this environment variable:


```bash
export GOOGLE_APPLICATION_CREDENTIALS="[PATH]"
```

Note: you may skip this step if you are running your Go application on GCP and it has access to Stackdriver Monitoring API.

## Usage

```go
import "github.com/plutov/gosd"

func main() {
    // This goroutine will send stats on your behalf
    go gosd.Run(gosd.Config{
		ProjectID: "PROJECT_ID",
		Logger:    os.Stdout,
		Labels:    map[string]string{"app": "my-web-app"},
	})
}
```

### Stackdriver Metrics

In Stackdriver go to Resources -> Metrics Explorer and find these metrics:

- custom.googleapis.com/gosd/goroutines
- custom.googleapis.com/gosd/cgocalls
- custom.googleapis.com/gosd/mstats/*

### Memory Consumption

**gosd** adds around 5Mi to your program memory usage.

### Stackdriver Dashboard

![dashboard.png](https://raw.githubusercontent.com/plutov/gosd/master/dashboard.png)