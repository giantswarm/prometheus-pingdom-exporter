package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server [username] [password] [api-key]",
		Short: "Start the HTTP server",
		Run:   serverRun,
	}

	waitSeconds int
	port        int

	pingdomUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pingdom_up",
		Help: "Whether the last pingdom scrape was successfull (1: up, 0: down)",
	})

	pingdomCheckStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_check_status",
		Help: "The current status of the check (0: up, 1: unconfirmed_down, 2: down, -1: paused, -2: unknown)",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})

	pingdomCheckResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_check_response_time",
		Help: "The response time of last test in milliseconds",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})
)

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&waitSeconds, "wait", 10, "time (in seconds) between accessing the Pingdom  API")
	serverCmd.Flags().IntVar(&port, "port", 8000, "port to listen on")

	prometheus.MustRegister(pingdomUp)
	prometheus.MustRegister(pingdomCheckStatus)
	prometheus.MustRegister(pingdomCheckResponseTime)
}

func sleep() {
	time.Sleep(time.Second * time.Duration(waitSeconds))
}

func serverRun(cmd *cobra.Command, args []string) {
	var client *pingdom.Client
	flag.Parse()

	if len(cmd.Flags().Args()) == 3 {
		client = pingdom.NewClient(
			flag.Arg(1),
			flag.Arg(2),
			flag.Arg(3),
		)
	} else if len(cmd.Flags().Args()) == 4 {
		client = pingdom.NewMultiUserClient(
			flag.Arg(1),
			flag.Arg(2),
			flag.Arg(3),
			flag.Arg(4),
		)
	} else {
		cmd.Help()
		os.Exit(1)
	}

	go func() {
		var oldCheckMetrics map[int]prometheus.Labels
		for {
			params := map[string]string{
				"include_tags": "true",
			}
			checks, err := client.Checks.List(params)
			if err != nil {
				log.Println("Error getting checks ", err)
				pingdomUp.Set(0)

				sleep()
				continue
			}
			pingdomUp.Set(1)

			checkMetrics := make(map[int]prometheus.Labels)
			for _, check := range checks {
				id := strconv.Itoa(check.ID)

				var status float64
				switch check.Status {
				case "unknown":
					status = -2
				case "paused":
					status = -1
				case "up":
					status = 0
				case "unconfirmed_down":
					status = 1
				case "down":
					status = 2
				default:
					status = 100
				}

				resolution := strconv.Itoa(check.Resolution)

				paused := strconv.FormatBool(check.Paused)
				// Pingdom library doesn't report paused correctly,
				// so calculate it off the status.
				if check.Status == "paused" {
					paused = "true"
				}

				var tagsRaw []string
				for _, tag := range check.Tags {
					tagsRaw = append(tagsRaw, tag.Name)
				}
				tags := strings.Join(tagsRaw, ",")

				labels := map[string]string{
					"id": id,
					"name": check.Name,
					"hostname": check.Hostname,
					"resolution": resolution,
					"paused": paused,
					"tags": tags,
				}

				pingdomCheckStatus.With(labels).Set(status)
				pingdomCheckResponseTime.With(labels).Set(float64(check.LastResponseTime))

				checkMetrics[check.ID] = labels
			}

			for id, oldLabels := range oldCheckMetrics {
				if labels, found := checkMetrics[id]; !found || !reflect.DeepEqual(oldLabels, labels) {
					pingdomCheckStatus.Delete(oldLabels)
					pingdomCheckResponseTime.Delete(oldLabels)
				}
			}

			oldCheckMetrics = checkMetrics

			sleep()
		}
	}()

	go func() {
		intChan := make(chan os.Signal)
		termChan := make(chan os.Signal)

		signal.Notify(intChan, syscall.SIGINT)
		signal.Notify(termChan, syscall.SIGTERM)

		select {
		case <-intChan:
			log.Print("Received SIGINT, exiting")
			os.Exit(0)
		case <-termChan:
			log.Print("Received SIGTERM, exiting")
			os.Exit(0)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	http.Handle("/metrics", prometheus.Handler())

	log.Print("Listening on port ", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
