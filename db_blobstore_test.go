package wattwerks

import (
	"appengine/aetest"
        "bytes"
        "fmt"
        "io"
        "mime/multipart"
        "net/http"
        "net/http/httptest"
        "os"
        "path/filepath"
	"testing"
        "net/url"
)

func TestBlob(t *testing.T) {

        inst, err := aetest.NewInstance(nil)
        if err != nil {
                t.Fatal(err)
        }
        defer inst.Close()

        c1, err := aetest.NewContext(nil)
        if err != nil {
                t.Fatal(err)
        }
        defer c1.Close()
        bs := &BS{ctx: c1}

        var url *url.URL
        if url, err = bs.UploadURL("admin/goods/upload"); url == nil || err != nil {
                t.Errorf("%v", err.Error())
        }

        path, _ := os.Getwd()
        path += "/products_no_header.csv"
        extraParams := map[string]string {
                "title": "products_no_header.csv",
                "author": "doofus chan",
                "description": "csv file",
        }
        req, err := newInstRequest(inst, url.String(), extraParams, "file", path)
        if err != nil {
                t.Fatal(err)
        }

        //subs := make([]byte, 1000)
        //_,_ = req.Body.Read(subs)
        //t.Errorf("Subs %v", string(subs))

        //if err != nil {
        //        t.Fatal(err)
        //}
        w := httptest.NewRecorder()
        handleGoodsCreateCc(w, req)

        //t.Errorf("%d ~ %s", w.Code, w.Body.String())

}


func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer file.Close()

        body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        part, err := writer.CreateFormFile(paramName, filepath.Base(path))
        if err != nil {
                return nil, err
        }

        _, err = io.Copy(part, file)

        for k,v := range params {
                _ = writer.WriteField(k, v)
        }
        err = writer.Close()
        if err != nil {
                return nil, err
        }
        req, err := http.NewRequest("POST", uri, body)
        req.Header.Set("Content-Type", writer.FormDataContentType())
        return req, err

}

func newInstRequest(inst aetest.Instance, uri string, params map[string]string, paramName, path string) (*http.Request, error) {

        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer file.Close()

        /*body := &bytes.Buffer{}
        writer := multipart.NewWriter(body)
        part, err := writer.CreateFormFile(paramName, filepath.Base(path))
        if err != nil {
                return nil, err
        }

        _, err = io.Copy(part, file)

        for k,v := range params {
                _ = writer.WriteField(k, v)
        }
        err = writer.Close()
        if err != nil {
                return nil, err
        }
        req, err := inst.NewRequest("POST", uri, body)
        req.Header.Set("Content-Type", writer.FormDataContentType())
        return req, err*/
        fmt.Println(filepath.Base(path))

        bodyReader, bodyWriter := io.Pipe()
        multiWriter := multipart.NewWriter(bodyWriter)
        //defer multiWriter.Close()
        //errChan := make(chan error, 1)
        go func() {
                defer bodyWriter.Close()
                //defer file.Close()
                _ = multiWriter.SetBoundary("doofus")
                part, err := multiWriter.CreateFormFile(paramName, filepath.Base(path))
                if err != nil {
                        return
                }
                if _, err := io.Copy(part, file); err != nil {
                        return
                }
                for k,v := range params {
                        if err := multiWriter.WriteField(k,v); err != nil {
                                return
                        }
                }
                if err := multiWriter.Close(); err != nil {
                        return
                }
        }()
        req, err := inst.NewRequest("POST", uri, bodyReader)
        req.Header.Set("Content-Type", multiWriter.FormDataContentType())
        return req, err
}
