package wattwerks

import (
        _"fmt"
        "math/rand"
        "reflect"
        "testing"

        "appengine/aetest"
        _"google.golang.org/appengine/memcache"
)

var (
        geterr = "Get error: pointer reqd"

)

func newGood() (*Good, *GoodDeets) {

        pd1 := GoodDeets{ DescDetails: "", Price: 100.00, Tax: 0.155, Stock: 12, Related:[]int64{}, Prices: []float64{}, Volumes: []int{}, ParameterNames: []string{}, ParameterValues: []string{}, Features: []string{}, Items: []string{}, UrlImgs1: "" , UrlImgs2: "", UrlImgs3: "", UrlFile: ""}
        return &Good{ Id: rand.Int63n(50), Code: "Code", Category: "Inverter", Subcategory: "Micro", Brand: "SMA", Desc: "250W Microinverter", Price: 100.00, Url: "", Urlimg: "", Featured: true, Hidden: false, Deets: pd1 }, &pd1

}

func TestAddInterface(t *testing.T) {

        c, err := aetest.NewContext(nil)
        if err != nil {
                t.Fatal(err)
        }
        defer c.Close()

        adsc := &DS{ctx: c}
        g, gd := newGood()
        if _, err := adsc.Add(g, g.Id); err != nil {
                t.Errorf("Good creation error: %v", err)
        }

        g1 := &Good{Id: g.Id, Deets:*gd } // Deets needs to be set with non nil slices because datastore will return a nil slice
        if err := adsc.Get(g1);  err != nil || !reflect.DeepEqual(g, g1)  {
                if err != nil {
                        t.Errorf("Good get error: %v", err)
                } else {
                        t.Errorf("Expect: %v, Got: %v, %v %v", g, g1, reflect.ValueOf(g.Deets.Prices).IsNil(), reflect.ValueOf(g1.Deets.Prices).IsNil())
                }
        }

        var g2 Good
        if err := adsc.Get(g2);  err.(DSErr).What != geterr {
                t.Errorf("Wanted: %v, Got: %v", geterr, err.(DSErr).What)
        }

        if err := adsc.AddwParent(g, gd, 1e7); err != nil {
                t.Errorf("Wanted: %v, Got: %v", nil, err.(DSErr).What)
        }

        g1 = &Good{Id: g.Id, Deets: *gd}
        if err := adsc.Get(g1);  err != nil || !reflect.DeepEqual(g, g1)  {
                if err != nil {
                        t.Errorf("Good get error: %v", err)
                } else {
                        t.Errorf("Expect: %v, Got: %v, %v %v", g, g1, reflect.ValueOf(g.Deets.Prices).IsNil(), reflect.ValueOf(g1.Deets.Prices).IsNil())
                }
        }

}
