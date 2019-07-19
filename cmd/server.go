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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server -username [username] -password [password] -api_key [api-key]",
		Short: "Start the HTTP server",
		Run:   serverRun,
	}

	waitSeconds int
	port        int

	// Pingdom authentication
	username     string
	password     string
	apiKey       string
	accountEmail string // Optional: For multi user client

	detailedTags string

	pingdomUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pingdom_up",
		Help: "Whether the last pingdom scrape was successfull (1: up, 0: down)",
	})

	labels             = []string{"id", "name", "hostname", "resolution", "paused", "tags", "region", "country", "city"}
	pingdomCheckStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_check_status",
		Help: "The current status of the check (0: up, 1: unconfirmed_down, 2: down, -1: paused, -2: unknown)",
	}, labels)

	pingdomCheckResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_check_response_time",
		Help: "The response time of last test in milliseconds",
	}, labels)
)

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&waitSeconds, "wait", 10, "time (in seconds) between accessing the Pingdom  API")
	serverCmd.Flags().IntVar(&port, "port", 8000, "port to listen on")
	serverCmd.Flags().StringVar(&username, "username", "", "Pingdom username")
	serverCmd.Flags().StringVar(&password, "password", "", "Pingdom password")
	serverCmd.Flags().StringVar(&apiKey, "api_key", "", "Pingdom api-key")
	serverCmd.Flags().StringVar(&accountEmail, "account_email", "", "Pingdom account_email (Optional, for multi user client)")
	serverCmd.Flags().StringVar(&detailedTags, "detailed_tags", "", "Comma separated list of tags to match checks to enhance (probe infos)")

	prometheus.MustRegister(pingdomUp)
	prometheus.MustRegister(pingdomCheckStatus)
	prometheus.MustRegister(pingdomCheckResponseTime)
}

func sleep() {
	time.Sleep(time.Second * time.Duration(waitSeconds))
}

func serverRun(cmd *cobra.Command, args []string) {
	flag.Parse()

	go handleSignals()

	var client *pingdom.Client
	if username == "" || password == "" || apiKey == "" {
		fmt.Fprintf(os.Stderr, "username, password and api_key are mandatory")
		cmd.Help()
		os.Exit(1)
	}
	if accountEmail == "" {
		client = pingdom.NewClient(username, password, apiKey)
	} else {
		client = pingdom.NewMultiUserClient(username, password, apiKey, accountEmail)
	}

	tags := make(map[string]struct{})
	for _, t := range strings.Split(detailedTags, ",") {
		tags[t] = struct{}{}
	}

	go func() {
		for {
			if err := fetchMetrics(client, tags); err != nil {
				log.Print(err)
			}
			sleep()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "running")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Print("Listening on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func fetchMetrics(client *pingdom.Client, tags map[string]struct{}) error {
	probes := make(map[int]pingdom.ProbeResponse)
	if len(tags) > 0 {
		pr, err := client.Probes.List()
		if err != nil {
			pingdomUp.Set(0)
			return fmt.Errorf("Error getting probes: %v", err)
		}
		for _, p := range pr {
			probes[p.ID] = p
		}
	}

	checks, err := client.Checks.List(map[string]string{"include_tags": "true"})
	if err != nil {
		pingdomUp.Set(0)
		return fmt.Errorf("Error getting checks ", err)
	}
	pingdomUp.Set(1)

	for _, check := range checks {
		var details map[string]string
		for _, tag := range check.Tags {
			if _, ok := tags[tag.Name]; !ok {
				continue
			}
			check, details, err = checkDetails(client, probes, check)
			if err != nil {
				log.Printf("Failed to get check %v details: %v", err)
			}
			break
		}
		exportCheck(check, details)
	}
	return nil
}

func checkDetails(client *pingdom.Client, probes map[int]pingdom.ProbeResponse, check pingdom.CheckResponse) (pingdom.CheckResponse, map[string]string, error) {
	results, err := client.Checks.Results(check.ID, map[string]string{"limit": "1"})
	if err != nil {
		return check, nil, fmt.Errorf("Error getting check %v results: %v", check.ID, err)
	}
	if len(results.Results) == 0 {
		return check, nil, fmt.Errorf("No results found for check %v", check.ID)
	}
	last := results.Results[0]
	probe, ok := probes[last.ProbeID]
	if !ok {
		return check, nil, fmt.Errorf("Probe %v for check %v not found", last.ProbeID, check.ID)
	}

	// Take last result and override check data in case there was a check between
	// List and Results.
	check.LastResponseTime = int64(last.ResponseTime)
	check.Status = last.Status

	return check, map[string]string{"region": probe.Region, "country": probe.CountryISO, "city": probe.City}, nil
}

func parseStatus(s string) (status float64) {
	switch s {
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
	return
}

var publishedLabels = make(map[int][]string)

func exportCheck(check pingdom.CheckResponse, details map[string]string) {
	var (
		id         = strconv.Itoa(check.ID)
		status     = parseStatus(check.Status)
		resolution = strconv.Itoa(check.Resolution)
		paused     = strconv.FormatBool(check.Paused)
		region     = ""
		country    = ""
		city       = ""
	)
	// Pingdom library doesn't report paused correctly,
	// so calculate it off the status.
	if check.Status == "paused" {
		paused = "true"
	}
	if details != nil {
		region = details["region"]
		country = details["country"]
		city = details["city"]
	}

	var tagsRaw []string
	for _, tag := range check.Tags {
		tagsRaw = append(tagsRaw, tag.Name)
	}
	tags := strings.Join(tagsRaw, ",")

	values := []string{
		id,
		check.Name,
		check.Hostname,
		resolution,
		paused,
		tags,
		region,
		country,
		city,
	}
	if old, ok := publishedLabels[check.ID]; ok && !reflect.DeepEqual(old, values) {
		pingdomCheckStatus.DeleteLabelValues(old...)
		pingdomCheckResponseTime.DeleteLabelValues(old...)
	}
	publishedLabels[check.ID] = values

	pingdomCheckStatus.WithLabelValues(values...).Set(status)
	pingdomCheckResponseTime.WithLabelValues(values...).Set(float64(check.LastResponseTime))
}

func handleSignals() {
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
}
