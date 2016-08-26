package wattwerks

import (
        "appengine"
	"appengine/blobstore"
	"appengine/datastore"
        "bytes"
        "github.com/tortuoise/csv"
        "fmt"
        "log"
        "net/http"
        "time"
        "net/url"
)

type BS struct {

        ctx appengine.Context

}

type BSErr struct {
        When time.Time
        What string
}

func (bs *BS) UploadURL(successPath string) (*url.URL, error) {

        if bs.ctx == nil {
                return nil, BSErr{When: time.Now(), What: "Context is nil"}
        }

        upURL, err := blobstore.UploadURL(bs.ctx, successPath, nil)
        if err != nil {
                return nil, BSErr{When: time.Now(), What:"UploadURL: " + err.Error()}
        }
        return upURL, nil
}

func (bs *BS) ParseUpload(r *http.Request, s string) ([]byte, error) {

        if bs.ctx == nil {
                return nil, BSErr{When: time.Now(), What: "Context is nil",}
        }

        blobs, _, err := blobstore.ParseUpload(r)
        if err != nil {
                return nil, BSErr{When: time.Now(), What: "ParseUpload: " + err.Error()}
        }

        var ok bool
        file := make([]*blobstore.BlobInfo,1)
        if file,ok = blobs[s]; !ok || len(file) == 0{
                return nil, BSErr{When: time.Now(), What: "Couldn't locate blob" + err.Error()}
        }

        rdr := blobstore.NewReader(bs.ctx, file[0].BlobKey)
        ubs := make([]byte, file[0].Size)
        if n, err := rdr.Read(ubs); n == 0 || err != nil {
                return nil, BSErr{When: time.Now(), What: "Couldn't read blob" + err.Error()}
        }

        return ubs, nil

}

func (bs *BS) UnmarshalCSV(data []byte, n int) error {

        dec := csv.NewDecoder(bytes.NewReader(data))
        sc := make(chan string)
        done := make(chan bool)
        for i := 0; i < n; i++ {
                go func(id string) {
                        sc <- fmt.Sprintf("%s - running", id)
                        var good Good
                        var na int

                        if err := dec.DeepUnmarshalCSV(&good); err != nil {
                                done <- true
                                return
                        }
                        err := datastore.RunInTransaction(bs.ctx, func(c appengine.Context) error {

                                na ++
                                return persist(c, sc, good, na)

                        }, nil)
                        if err != nil {
                                sc <- fmt.Sprintf("%s: failure %v\n", id, err)
                        } else {
                                sc <- fmt.Sprintf("%s: success %v\n", id)
                        }

                        done <- true
                }(fmt.Sprintf("#%04d", 1 << uint(i)))
        }
        for nDone := 0; nDone < n; {
                select {
                case s := <-sc:
                        log.Println(s)
                case <-done:
                        nDone++
                }
        }
        return nil
}


func persist(c appengine.Context, sc chan string, good Good, na int) error {

        ds := &DS{ctx:c}
        if err := ds.AddwParent(&good, &good.Deets, 1e7); err != nil {
                return BSErr{When: time.Now(), What: "Couldn't persist" + err.Error()}
        }
        return nil
}

func (bs *BS) UnmarshalJSON() {


}


func (bse BSErr) Error() string {

        return fmt.Sprintf("%v, %v", bse.When, bse.What)

}
