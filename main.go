package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func init() {
	flag.String("port", "80", "aws endpoint to use")
	flag.Bool("debug", false, "sqs url to poll")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigName("config")

	viper.AutomaticEnv()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}

func main() {
	debug := viper.GetBool("debug")
	port := viper.GetString("port")
	r := NewRouter()
	log.Warnf("GoAws listening on: 0.0.0.0:%s debug:%v", port, debug)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
