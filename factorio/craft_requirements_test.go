package factorio

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestSingleItemComponents(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	dbString, err := ioutil.ReadFile("testdata/itemdb.json")
	if err != nil {
		t.Error("Failed to read DB file:", err.Error())
		t.Fail()
	}
	db, dbErr := LoadJsonToDb(string(dbString))
	if dbErr != nil {
		t.Error("Failed to parse DB file:", dbErr.Error())
		t.Fail()
	}

	testItemId := "boiler"
	var report = CraftComponents{Items: make(map[string]int)}
	computeErr := db.ComputeCraftComponents(testItemId, 1, &report)
	if computeErr != nil {
		t.Error("Error computing components for a", testItemId, ":", err.Error())
		t.Fail()
	}
	if 5 != len(report.Items) {
		t.Errorf("Incorrect component count. Expected %d, got %d",
			5,
			len(report.Items))
	}
	if 4 != report.Items["pipe"] {
		t.Errorf("Incorrect pipe count. Expected %d, got %d",
			4,
			report.Items["pipe"])
	}
}

func TestRecursiveComputation(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	dbString, err := ioutil.ReadFile("testdata/itemdb.json")
	if err != nil {
		t.Error("Failed to read DB file:", err.Error())
		t.Fail()
	}
	db, dbErr := LoadJsonToDb(string(dbString))
	if dbErr != nil {
		t.Error("Failed to parse DB file:", dbErr.Error())
		t.Fail()
	}

	testItemId := "boiler"
	var report = CraftComponents{Items: make(map[string]int)}
	computeErr := db.ComputeCraftComponents(testItemId, 1, &report)
	if computeErr != nil {
		t.Error("Error computing components for a", testItemId, ":", err.Error())
		t.Fail()
	}
	if 5 != len(report.Items) {
		t.Errorf("Incorrect component count. Expected %d, got %d",
			5,
			len(report.Items))
	}
	if 5 != report.Items["stone"] {
		t.Errorf("Incorrect stone count. Expected %d, got %d",
			5,
			report.Items["stone"])
	}

	if 4 != report.Items["iron-plate"] {
		t.Errorf("Incorrect iron-plate count. Expected %d, got %d",
			4,
			report.Items["iron-plate"])
	}
	if 4 != report.Items["iron-ore"] {
		t.Errorf("Incorrect iron-ore count. Expected %d, got %d",
			4,
			report.Items["iron-ore"])
	}
}
