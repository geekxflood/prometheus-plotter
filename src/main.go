package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("prometheus_url", "http://localhost:9090")
	viper.SetDefault("start_time", time.Now().Add(-24*time.Hour).Unix())
	viper.SetDefault("end_time", time.Now().Unix())
	viper.SetDefault("start_time_str", "")
	viper.SetDefault("end_time_str", "")
	viper.SetDefault("query", "")
}

func main() {
	var configFile string
	var query string
	var err error

	rootCmd := &cobra.Command{
		Use:   "prometheus-query-plotter",
		Short: "A tool for plotting data from a Prometheus server",
		Run: func(cmd *cobra.Command, args []string) {
			// Load configuration from file
			if configFile != "" {
				viper.SetConfigFile(configFile)
				err := viper.ReadInConfig()
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					os.Exit(1)
				}
			}

			// Get Prometheus URL from configuration file or command-line option
			url := viper.GetString("prometheus_url")
			if url == "" {
				fmt.Println("Error: Prometheus URL not specified")
				os.Exit(1)
			}

			// Get query from configuration file or user prompt
			if query == "" {
				query = viper.GetString("query")
				if query == "" {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Enter query: ")
					query, _ = reader.ReadString('\n')
					query = strings.TrimRight(query, "\n")
				}
			}

			// Get time range from configuration file or command-line options
			startTime := viper.GetInt64("start_time")
			endTime := viper.GetInt64("end_time")
			if startTime == 0 || endTime == 0 {
				startTimeStr := viper.GetString("start_time_str")
				endTimeStr := viper.GetString("end_time_str")
				if startTimeStr == "" || endTimeStr == "" {
					fmt.Println("Error: Time range not specified")
					os.Exit(1)
				}

				startTime, err = parseTime(startTimeStr)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					os.Exit(1)
				}

				endTime, err = parseTime(endTimeStr)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					os.Exit(1)
				}
			}

			// Retrieve data from Prometheus
			data, err := getData(url, query, startTime, endTime)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			// Extract the relevant data points
			values, err := extractData(data)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			// Create a plot of the data
			output := viper.GetString("output")
			if output == "" {
				output = "output.png"
			}
			err = createPlot(values, "Prometheus Query", "Time", "Requests", output)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}

			fmt.Println("Done!")
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "Prometheus query")

	viper.SetDefault("output", "output.png")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("prometheus_query")
	viper.AutomaticEnv()

	err = rootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
