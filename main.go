package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"runtime/debug"

	collectd "github.com/paulhammond/gocollectd"
	"gopkg.in/yaml.v2"
)

func parseCommandLineParams() {
	flag.StringVar(&configPath, "c", "./config.yml", "Path to config.yml")
	flag.Parse()
}

func initConfigs() {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("error reading config", err)
	}
}

func initRuntime() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	log.Printf("Init runtime to use %d CPUs and %d threads", numCPU, config.System.MaxThreads)
	debug.SetMaxThreads(config.System.MaxThreads)
}

func main() {
	fmt.Printf("Version:    [%s]\nBuild:      [%s]\nBuild Date: [%s]\n", version, build, buildDate)
	parseCommandLineParams()
	initConfigs()
	initRuntime()
	connectClickDB()

	c := make(chan collectd.Packet)
	log.Println("Starting listner on " + config.System.ListenOn)
	go collectd.Listen(config.System.ListenOn, c)
	for {
		packet := <-c
		// do something with the packet
		processCollectDPacket(packet)
	}
}

func processCollectDPacket(packet collectd.Packet) {

	hostname := packet.Hostname
	pkgtime := packet.Time()   // A go time value
	pkgplugin := packet.Plugin // "Load"
	values, err := packet.ValueNumbers()
	if err != nil {
		log.Fatalln("Wrong packet from collectd", err)
	}
	names := packet.ValueNames()

	var insData []cdElementData

	for i, dat := range names {
		var tmpIns cdElementData
		tmpIns.EventDateTime = pkgtime
		tmpIns.Hostname = hostname
		tmpIns.Plugin = pkgplugin
		tmpIns.ParamName = dat
		tmpIns.ParamValue = values[i].Float64()
		insData = append(insData, tmpIns)
	}

	err = insertCollectDToCH(insData)
	if err != nil {
		log.Fatalln("Wrong inserting data to Clickhouse", err)
	}
}
