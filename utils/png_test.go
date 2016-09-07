package utils

import (
        "path/filepath"
	_"fmt"
        "regexp"
        "os"
        "strings"
        "testing"
)

var JPG = regexp.MustCompile("^.[j|J][p|P][e|E]?[g|G]$")
var PNG = regexp.MustCompile("^.[p|P][n|N][g|G]$")
var DB = regexp.MustCompile("\\s+\\w+\\s+from")
var RM = regexp.MustCompile(".+\\.png.+\\.jpg")
var URL = regexp.MustCompile("/paste/v1/([^/]+)")

// TestConvert walks filepath and saves a jpg format image from all png format images 
func TestConvert(t *testing.T) {


        err := filepath.Walk("../img/goods/med/", func(path string, info os.FileInfo, err error) error {

                if err != nil {
                        t.Logf("%v", err)
                }
                if !info.IsDir() {
                        rdr, err := os.Open(path)
                        if err != nil {
                                //t.Logf("%v", err)
                                return err
                        }
                        defer rdr.Close()
                        sxt := filepath.Ext(rdr.Name())
                        if ok := PNG.Find([]byte(sxt)); len(ok) == 0 {
                                return err
                        }
                        ss := strings.Split(info.Name(), ".")
                        wtr, err := os.Create("../img/goods/med/" + ss[0] + ".jpg")
                        if err != nil {
                                t.Logf("%v", err)
                        }
                        defer wtr.Close()

                        //err = convertJPG2PNG(wtr, rdr)
                        err = convertPNG2JPG(wtr, rdr)
                        if err != nil {
                                t.Logf("%v", err)
                        }
                        return nil
                } else {
                        t.Logf("Skipped directory: %s", info.Name())
                }
                return nil
        })
        if err != nil {
                t.Errorf("%v", err)
        }

}


func TestRegexp(t *testing.T) {

        s := "select uname as name weight AS wgt from table"
        f := DB.Find([]byte(s))
        t.Logf("%v", string(f))

        s = ".JPG"
        f = JPG.Find([]byte(s))
        t.Logf("%v", string(f))

        s = "xy_adf.pngadfdssf.jpg"
        f = RM.Find([]byte(s))
        t.Logf("%v", string(f))

        s = "/paste/v1/aSdf3"
        g := URL.FindStringSubmatch(s)
        t.Logf("%v", g)
}
