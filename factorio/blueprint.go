package factorio

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Blueprint struct {
	Details BlueprintDetails `json:"blueprint"`
}

type BlueprintDetails struct {
	Icons    []Icon   `json:"icons"`
	Entities []Entity `json:"entities"`
	Item     string   `json:"item"`
	Label    string   `json:"label"`
	Version  int      `json:"version"`
}
type Icon struct {
	Signal IconSignal `json:"signal"`
	Index  int        `json:"index"`
}
type IconSignal struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
type Entity struct {
	Number    int            `json:"entity_number"`
	Position  EntityPosition `json:"position"`
	Name      string         `json:"name"`
	Direction int            `json:"direction"`
	Type      string         `json:"type"`
}
type EntityPosition struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type BlueprintBookEnvelope struct {
	BlueprintBook BlueprintBook `json:"blueprint_book"`
}

type BlueprintBook struct {
	Blueprints  []Blueprint `json:"blueprints"`
	Item        string      `json:"item"`
	Label       string      `json:"label"`
	ActiveIndex int         `json:"active_index"`
	Version     int         `json:"version"`
}

func ParseBlueprintString(bp string) (blueprint *Blueprint, err error) {

	actualBase64String := strings.TrimPrefix(bp, "0")
	data, decodeErr := base64.StdEncoding.DecodeString(actualBase64String)
	if decodeErr != nil {
		log.Errorf("Failed to decode the Base64 string:", decodeErr.Error())
		return nil, decodeErr
	}

	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	defer r.Close()
	if err != nil {
		log.Errorln("Failed to construct zlib reader:", err)
		return nil, err
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)
	bpString := buffer.String()

	// Deserialize the JSON data
	blueprint = new(Blueprint)
	err = json.Unmarshal([]byte(bpString), blueprint)
	if err != nil {
		log.Errorln("Failed to unmarshal BP:", err)
	}
	return
}

func ParseBlueprintBookString(bp string) (bpBook *BlueprintBook, err error) {
	data, decodeErr := base64.StdEncoding.DecodeString(bp)
	if decodeErr != nil {
		log.Errorf("Failed to decode the Base64 string: %s", decodeErr)
		return nil, decodeErr
	}

	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	defer r.Close()
	if err != nil {
		log.Errorln("Failed to construct zlib reader:", err)
		return nil, err
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r)
	bpString := buffer.String()

	// Deserialize the JSON data
	bpBookEnvelope := new(BlueprintBookEnvelope)
	err = json.Unmarshal([]byte(bpString), bpBookEnvelope)
	if err != nil {
		log.Errorln("Failed to unmarshal BP:", err)
		return nil, err
	}
	return &bpBookEnvelope.BlueprintBook, nil
}
