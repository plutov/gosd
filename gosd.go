package gosd

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Push interval (1 min)
const intervalSeconds = 60

// map of all metrics with name as a key and func go get a data point
var metrics = map[string]func() *monitoringpb.Point{
	"custom.googleapis.com/gosd/goroutime":           getNumGoroutineDataPoint,
	"custom.googleapis.com/gosd/cgocall":             getNumCgoCallDataPoint,
	"custom.googleapis.com/gosd/mstats/alloc":        getMemStatsAllocDataPoint,
	"custom.googleapis.com/gosd/mstats/totalalloc":   getMemStatsTotalAllocDataPoint,
	"custom.googleapis.com/gosd/mstats/sys":          getMemStatsSysDataPoint,
	"custom.googleapis.com/gosd/mstats/mallocs":      getMemStatsMallocsDataPoint,
	"custom.googleapis.com/gosd/mstats/frees":        getMemStatsFreesDataPoint,
	"custom.googleapis.com/gosd/mstats/pausetotalns": getMemStatsPauseTotalNsDataPoint,
	"custom.googleapis.com/gosd/mstats/numgc":        getMemStatsNumGCDataPoint,
	"custom.googleapis.com/gosd/mstats/GCSys":        getMemStatsGCSysDataPoint,
	"custom.googleapis.com/gosd/mstats/lookups":      getMemStatsLookupsDataPoint,
	"custom.googleapis.com/gosd/mstats/heapalloc":    getMemStatsHeapAllocDataPoint,
	"custom.googleapis.com/gosd/mstats/heapsys":      getMemStatsHeapSysDataPoint,
	"custom.googleapis.com/gosd/mstats/heapidle":     getMemStatsHeapIdleDataPoint,
	"custom.googleapis.com/gosd/mstats/heapinuse":    getMemStatsHeapInuseDataPoint,
	"custom.googleapis.com/gosd/mstats/HeapReleased": getMemStatsHeapReleasedDataPoint,
	"custom.googleapis.com/gosd/mstats/HeapObjects":  getMemStatsHeapObjectsDataPoint,
	"custom.googleapis.com/gosd/mstats/StackInuse":   getMemStatsStackInuseDataPoint,
	"custom.googleapis.com/gosd/mstats/StackSys":     getMemStatsStackSysDataPoint,
}

var rtm runtime.MemStats

// Run starts a goroutine which will automatically send Go runtime stats to stackdriver.
// projectID is your Google Cloud Platform project ID.
func Run(projectID string, logger io.Writer) {
	ctx := context.Background()

	// Creates a client
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		fmt.Fprintf(logger, "unable to create monitoring client: %s", err.Error())
		return
	}

	ticker := time.NewTicker(time.Second * intervalSeconds)
	for range ticker.C {
		runtime.ReadMemStats(&rtm)
		var timeSeries []*monitoringpb.TimeSeries

		for name, f := range metrics {
			timeSeries = append(timeSeries, &monitoringpb.TimeSeries{
				Metric: &metricpb.Metric{
					Type: name,
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
				},
				Points: []*monitoringpb.Point{
					f(),
				},
			})
		}
		// write all metrics at once
		if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name:       monitoring.MetricProjectPath(projectID),
			TimeSeries: timeSeries,
		}); err != nil {
			fmt.Fprintf(logger, "unable to write time series data: %s", err.Error())
		}
	}
}

// since all our metrics are int64, we use this helper to return a data point with value
func getInt64DataPoint(value int64) *monitoringpb.Point {
	return &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: value,
			},
		},
	}
}

func getNumGoroutineDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(runtime.NumGoroutine()))
}

func getNumCgoCallDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(runtime.NumCgoCall())
}

func getMemStatsAllocDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.Alloc))
}

func getMemStatsTotalAllocDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.TotalAlloc))
}

func getMemStatsSysDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.Sys))
}

func getMemStatsMallocsDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.Mallocs))
}

func getMemStatsFreesDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.Frees))
}

func getMemStatsPauseTotalNsDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.PauseTotalNs))
}

func getMemStatsNumGCDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.NumGC))
}

func getMemStatsLookupsDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.Lookups))
}

func getMemStatsHeapAllocDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapAlloc))
}

func getMemStatsHeapSysDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapSys))
}

func getMemStatsHeapIdleDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapIdle))
}

func getMemStatsHeapInuseDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapInuse))
}

func getMemStatsHeapReleasedDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapReleased))
}

func getMemStatsHeapObjectsDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.HeapObjects))
}

func getMemStatsStackInuseDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.StackInuse))
}

func getMemStatsStackSysDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.StackSys))
}

func getMemStatsGCSysDataPoint() *monitoringpb.Point {
	return getInt64DataPoint(int64(rtm.GCSys))
}
