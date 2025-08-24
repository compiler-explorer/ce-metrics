package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr = flag.String("listen-address", ":9100", "The address to listen on for HTTP requests.")
	
	mountsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ce_node_mounts_total",
			Help: "Total number of mount points on the node",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(mountsTotal)
}

func getMountInfo() error {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return err
	}
	defer file.Close()

	mountCounts := make(map[string]float64)
	var cefs1Count float64
	var cefs2Count float64
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		
		mountpoint := fields[1]
		fstype := fields[2]
		
		mountCounts[fstype]++
		mountCounts["all"]++
		
		// Count CEFS1 mounts at /efs/compiler-explorer/ with type squashfs
		if strings.HasPrefix(mountpoint, "/efs/compiler-explorer/") && fstype == "squashfs" {
			cefs1Count++
		}
		
		// Count CEFS2 mounts at /cefs/XX/* where XX is any 2-char hash prefix (these are from /efs/cefs-images)
		// Pattern: /cefs/XX/... where XX is exactly 2 characters
		if len(mountpoint) > 9 && mountpoint[8] == '/' && strings.HasPrefix(mountpoint, "/cefs/") {
			cefs2Count++
		}
	}
	
	for fstype, count := range mountCounts {
		mountsTotal.WithLabelValues(fstype).Set(count)
	}
	
	// Add CEFS1 and CEFS2 counts as separate types
	mountsTotal.WithLabelValues("cefs1").Set(cefs1Count)
	mountsTotal.WithLabelValues("cefs2").Set(cefs2Count)
	
	return scanner.Err()
}

func collectMetrics() {
	if err := getMountInfo(); err != nil {
		log.Printf("Error collecting mount info: %v", err)
	}
}

func main() {
	flag.Parse()
	
	collectMetrics()
	
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		collectMetrics()
		promhttp.Handler().ServeHTTP(w, r)
	})
	
	log.Printf("Starting CE Node Exporter on %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}