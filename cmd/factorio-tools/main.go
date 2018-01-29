//
// Copyright (C) 2018 Chris Cox
//

package main

import (
	"firststrikegames.com/factorio-tools/pkg/version"
	log "github.com/sirupsen/logrus"

)

func main() {
	log.WithFields(log.Fields{"service": "factorio-tools"}).Info("Initializing with version", version.VERSION)


}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello, world!")
}

