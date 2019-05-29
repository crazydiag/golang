package main

import (
//  "fmt"
  "os/exec"
  "regexp"
  "strings"
  "net/http"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "log"
  "time"
)

var (

  gauge = prometheus.NewGaugeVec(
     prometheus.GaugeOpts{
        Namespace: "session",
        Name:      "user",
        Help:      "This is my gauge",
     },
   []string{"terminal"},
   )
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(gauge)
}


func main() {

  http.Handle("/metrics", promhttp.Handler())
  r, _ := regexp.Compile("pts")
  rr, _ := regexp.Compile("tty")



  go func() {
     for {
         gauge.With(prometheus.Labels{"terminal":"pts"}).Set(float64(0))
         gauge.With(prometheus.Labels{"terminal":"tty"}).Set(float64(0))
        out, err := exec.Command("who").Output()
        if err != nil {
                log.Fatal(err)
        }
        s := strings.Split(string(out), "\n")
        c := regexp.MustCompile("[^\\s]+")
        for i := 0; i < len(s); i++ {
            d := c.FindAllString(s[i], -1)
//             fmt.Println(d)
            if len(d) > 0 {
            if r.MatchString(string(d[1])) {
            gauge.With(prometheus.Labels{"terminal":"pts"}).Inc()
            }
            if rr.MatchString(string(d[1])) {
            gauge.With(prometheus.Labels{"terminal":"tty"}).Inc()
            }

          }
        }
            time.Sleep(10 * time.Second)
//          time.Sleep(100 * time.Millisecond)
  }
  }()

  log.Fatal(http.ListenAndServe(":9080", nil))
}
