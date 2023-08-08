package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Version        = "Unknown build"
	BuildTimestamp = "unknown build timestamp."
)

type CommonConfig struct {
	BuildTimeStamp   time.Time
	GitVersion       string
	ShowVersion      bool
	Verbose          bool
	WorkingDirectory string

	NumCPU             int
	LoadRepetitions    int
	LoadDuration       int
	LoadPause          int
	LoadInitialTimeOut int
}

func sleep(cfg *CommonConfig, timeout int) {
	if timeout == -1 {
		timeout = randomInt(10)
	}
	verboseInfo(cfg, "Sleeping for %d seconds.\n", timeout)
	time.Sleep(time.Duration(timeout) * time.Second)
}

func randomOrInt(value int) string {
	if value == -1 {
		return "random"
	}
	return strconv.Itoa(value)
}

func verboseInfo(cfg *CommonConfig, format string, a ...any) (n int, err error) {
	if cfg.Verbose {
		n, err = fmt.Printf(format, a...)
	}
	return n, err
}

type worker struct {
	ch chan string
	ID int
}

func work(cfg *CommonConfig) {
	verboseInfo(cfg, "\n")
	var workers []*worker
	var wg sync.WaitGroup

	for i := 0; i < cfg.NumCPU; i++ {
		w := &worker{make(chan string), i}
		workers = append(workers, w)
		wg.Add(1)
		verboseInfo(cfg, "Starting worker %d.\n", w.ID)

		go func() {
			for {
				select {
				case cmd := <-w.ch:
					if cmd == "quit" {
						verboseInfo(cfg, "Stopped worker %d.\n", w.ID)
						defer wg.Done()
						close(w.ch)
						return
					}
				default:
					var nothing = float64(2)
					for i := 1; i < 10000; i++ {
						nothing = float64(i)
						nothing = math.Sqrt(float64(nothing))
					}
				}
			}
		}()
	}

	sleep(cfg, cfg.LoadDuration)
	verboseInfo(cfg, "\n")

	for _, w := range workers {
		verboseInfo(cfg, "Stopping worker %d .\n", w.ID)
		w.ch <- "quit"
		//close(w)
	}
	wg.Wait()
}

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

func randomInt(limit int) int {
	return r1.Intn(limit)
}

func main() {

	btime, err := time.Parse("2006-01-02_15:04:05_MST", BuildTimestamp)
	if err != nil {
		btime = time.Now()
	}
	cfg := &CommonConfig{}
	cfg.BuildTimeStamp = btime
	cfg.GitVersion = Version
	cfg.WorkingDirectory, _ = os.Getwd()
	cfg.NumCPU = runtime.NumCPU()

	flag.BoolVar(&cfg.ShowVersion, "version", false, "Show version info.")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Get debug info.")
	flag.IntVar(&cfg.LoadDuration, "d", 120, "Duration of a load peak to be generated in seconds.  Use -1 for random duration.")
	flag.IntVar(&cfg.LoadRepetitions, "r", -1, "Number of times the peak will be repeated. Use -1 for random times.")
	flag.IntVar(&cfg.LoadPause, "p", -1, "Fixed pause in seconds between the generated peaks. Use -1 for random amount.")
	flag.IntVar(&cfg.LoadInitialTimeOut, "t", -1, "Fixed pause in seconds before the first peakt will be created. Use -1 for random amount.")
	flag.Parse()

	if cfg.ShowVersion || cfg.Verbose {
		v := cfg.GitVersion
		if strings.ToLower(v) == "unknown build" {
			v = "'Unknown build'"
		}

		fmt.Printf("Version %s built on %s using %s.\n", v, cfg.BuildTimeStamp.Format("02.01.2006"), runtime.Version())
		if cfg.ShowVersion {
			os.Exit(0)
		}
	}
	if cfg.LoadDuration == -1 {
		cfg.LoadDuration = randomInt(10) + 1
	}
	if cfg.LoadRepetitions == -1 {
		cfg.LoadRepetitions = randomInt(2) + 1
	}
	if cfg.LoadInitialTimeOut == -1 {
		cfg.LoadInitialTimeOut = randomInt(10) + 17
	}

	verboseInfo(cfg, "Working Directory: %s\n", cfg.WorkingDirectory)
	verboseInfo(cfg, "LoadDuration: %s\n", randomOrInt(cfg.LoadDuration))
	verboseInfo(cfg, "LoadRepetitions: %s\n", randomOrInt(cfg.LoadRepetitions))
	verboseInfo(cfg, "LoadInitialTimeOut: %s\n", randomOrInt(cfg.LoadInitialTimeOut))
	verboseInfo(cfg, "LoadPause: %s\n\n", randomOrInt(cfg.LoadPause))

	sleep(cfg, cfg.LoadInitialTimeOut)
	for peaks := 0; peaks < cfg.LoadRepetitions; peaks++ {
		verboseInfo(cfg, "Generating peak %d\n", peaks)
		work(cfg)
		if peaks < cfg.LoadRepetitions-1 {
			sleep(cfg, cfg.LoadPause)
		}
	}
}
