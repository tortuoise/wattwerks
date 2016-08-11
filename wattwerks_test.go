package wattwerks

import (
	_"fmt"
	_"encoding/json"
	"encoding/csv"
	"bytes"
	_"io"
	"io/ioutil"
	"log"
	"testing"

	_"appengine"
	_"appengine/aetest"
	_"appengine/datastore"
	_"appengine/user"
)

func TestCSV(t *testing.T) {
	raw, err := ioutil.ReadFile("products.csv")
	handle(err)
	r := csv.NewReader(bytes.NewReader(raw))
	//for {
		records, err := r.ReadAll()
		if err != nil {
			log.Println(err)
		}
		if len(records) == 0 {
			t.Errorf("Boo")
		}
                gd := csvMarshal(records[1])
                if gd.Brand != records[1][4]{
                        t.Errorf("Boo")
                }
                if gd.Deets.ParameterValues[0] != records[1][45]{
                        t.Errorf("Boo")
                }
	//}
}

//func handle(err error) {
//	log.Println(err)
//}

