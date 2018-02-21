//
// Copyright (C) 2018 Chris Cox
//

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/Klaital/factorio-tools/factorio"
)

type appConfig struct {
	ItemDbPath	string
	AnalyzeItem	bool
}

func LoadConfigFromCli(config *appConfig) {
	flag.StringVar(&config.ItemDbPath,
			"item-path",
			"testdata/itemdb.json",
			"The JSON file containing the Item DB")
	flag.BoolVar(&config.AnalyzeItem,
			"analyze-item",
			false,
			"Analyze a named Items.")
	flag.Parse()
}

func main() {
	logger := log.WithFields(log.Fields{
		"func": "main",
	})
	log.SetLevel(log.DebugLevel)
	logger.Debugln("Initializing...")

	// Load configuration from commandline variables
	var config appConfig
	LoadConfigFromCli(&config)
	logger.Debugln("Config Loaded...")
	logger.Debugln("ItemDbPath:", config.ItemDbPath)

	dbString, err := ioutil.ReadFile(config.ItemDbPath)
	if err != nil {
		log.Errorf("failed to load data from file:", err)
		return
	}
	db, dbErr := factorio.LoadJsonToDb(string(dbString))
	if dbErr != nil {
		log.Errorln("failed to parse DB from JSON:", dbErr)
		return
	}

	log.Infoln("Loaded Database with", len(db.Data), "records")

	//
	// REAL WORK STARTS HERE
	//

	if config.AnalyzeItem {
		// Initialize the component breakdown report
		report := factorio.CraftComponents{Items:make(map[string]int)}

		for _, itemId := range flag.Args() {
			logger = logger.WithFields(log.Fields{
				"Item": itemId,
				"operation": "AnalyzeItem",
			})
			logger.Debugln("Analyzing Item:", itemId)

			err := db.ComputeCraftComponents(itemId, 1, &report)
			if err != nil {
				logger.Errorln("Failed to compute crafting requirements:", err.Error())
			} else {
				logger.Debugln(report)
				for itemId, qty := range report.Items {
					fmt.Println(itemId, " => ", qty)
				}
			}
		}
		return
	}
}

