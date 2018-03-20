package factorio

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestLoadItemDb(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	testString, err := ioutil.ReadFile("testdata/itemdb.json")
	if err != nil {
		t.Errorf("Failed to read Item DB from file:", err)
		t.Fail()
	}

	db, dbErr := LoadJsonToDb(string(testString))
	if dbErr != nil {
		t.Errorf("Failed to parse Item DB JSON string: %s", dbErr.Error())
		t.Fail()
	}

	if db == nil {
		t.Error("No ItemDb was generated")
		t.Fail()
	}

	if 0 == len(db.Data) {
		t.Error("No data loaded")
		t.Fail()
	}
	ironPlate, ironLoadOk := db.Data["iron-plate"]
	if !ironLoadOk {
		t.Error("Iron Plate was not found in the database")
	} else {
		if 1 != len(ironPlate.Recipes) {
			t.Error("Incorrect number of recipes for Iron Plate. Expected %d, got %d",
				1,
				len(ironPlate.Recipes))
		}
	}
}
