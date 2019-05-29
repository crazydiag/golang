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

func main() {

  http.Handle("/metrics", promhttp.Handler())

  prometheus.MustRegister(gauge)

  go func() {
     for {
        tty := 0 
        pts := 0 
        out, err := exec.Command("who").Output()
        if err != nil {
                log.Fatal(err)
        }
        r, _ := regexp.Compile("pts")
        rr, _ := regexp.Compile("tty")
//      fmt.Println(r.MatchString(string(out)))
//      fmt.Printf("%s", out)
        s := strings.Split(string(out), "\n")
        c := regexp.MustCompile("[^\\s]+")
        for i := 0; i < len(s); i++ {
            d := c.FindAllString(s[i], -1)
            if len(d) > 0 {
            if r.MatchString(string(d[1])) {
            pts++
            }
            if rr.MatchString(string(d[1])) {
            tty++
            }

//          fmt.Printf("%d           %s\n", i, s[i])
//            fmt.Printf("%s ------ %s\n", d[0], d[1])
          }
        }
        gauge.With(prometheus.Labels{"terminal":"pts"}).Set(float64(pts))
        gauge.With(prometheus.Labels{"terminal":"tty"}).Set(float64(tty))
        time.Sleep(100 * time.Second)
//        time.Sleep(1000 * time.Millisecond)
//        fmt.Printf("PTS --- %d ; TTY --- %d\n", pts , tty)
  }
  }()

  log.Fatal(http.ListenAndServe(":9080", nil))
}
