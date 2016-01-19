package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/koding/multiconfig"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Host  string `default:"127.0.0.1"`
	Port  int    `default:"12801"`
	River string `default:"http://127.0.0.1:12800/stat"`
}

var (
	masterBin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "master_bin",
		Help: "Master bin log file number",
	})

	masterPos = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "master_pos",
		Help: "Master bin log pos in file",
	})

	riverBin = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "river_bin",
		Help: "bin log file number that the river is reading",
	})

	riverPos = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "river_pos",
		Help: "bin log pos in file that the river is reading",
	})

	inserts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "insert_count",
		Help: "insert count",
	})

	updates = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_count",
		Help: "update count",
	})

	deletes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "delete_count",
		Help: "delete count",
	})
)

func handleBodyError(err error) {
	fmt.Println("error reading from native river stats output: ", err)
	masterBin.Set(1.0)
	riverBin.Set(0.0)
}

func readCounts(riverAddr string) {
	resp, err := http.Get(riverAddr)
	if err != nil {
		handleBodyError(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleBodyError(err)
		return
	}

	status := strings.Split(string(body), "\n")
	fmt.Println(status)
	// parse master bin
	re := regexp.MustCompile(`server.+\(.+\.(\d+)[^\d]+(\d+)\)`)
	strs := re.FindStringSubmatch(status[0])
	f, _ := strconv.ParseFloat(strs[1], 64)
	masterBin.Set(f)
	f, _ = strconv.ParseFloat(strs[2], 64)
	masterPos.Set(f)

	// parse river bin
	re = regexp.MustCompile(`read.+\(.+\.(\d+)[^\d]+(\d+)\)`)
	strs = re.FindStringSubmatch(status[1])
	f, _ = strconv.ParseFloat(strs[1], 64)
	riverBin.Set(f)
	f, _ = strconv.ParseFloat(strs[2], 64)
	riverPos.Set(f)

	re = regexp.MustCompile(`insert_num:(\d+)`)
	strs = re.FindStringSubmatch(status[2])
	f, _ = strconv.ParseFloat(strs[1], 64)
	inserts.Set(f)

	re = regexp.MustCompile(`update_num:(\d+)`)
	strs = re.FindStringSubmatch(status[3])
	f, _ = strconv.ParseFloat(strs[1], 64)
	updates.Set(f)

	re = regexp.MustCompile(`delete_num:(\d+)`)
	strs = re.FindStringSubmatch(status[4])
	f, _ = strconv.ParseFloat(strs[1], 64)
	deletes.Set(f)
}

func proxyHandler(h http.Handler, f func()) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f()
		h.ServeHTTP(w, r)
	})
}

func main() {
	mc := multiconfig.New()
	conf := new(Config)
	mc.MustLoad(conf)

	prometheus.MustRegister(masterBin)
	prometheus.MustRegister(masterPos)
	prometheus.MustRegister(riverBin)
	prometheus.MustRegister(riverPos)
	prometheus.MustRegister(inserts)
	prometheus.MustRegister(updates)
	prometheus.MustRegister(deletes)

	handler := proxyHandler(prometheus.Handler(), func() {
		readCounts(conf.River)
	})
	http.Handle("/metrics", handler)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), nil)
	if err != nil {
		panic(err)
	}
}
