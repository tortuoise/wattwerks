package wattwerks

import (
	"appengine/aetest"
	_ "bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
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
	extraParams := map[string]string{
		"title":       "products_no_header.csv",
		"author":      "doofus chan",
		"description": "csv file",
	}

	req, err := newInstRequest(inst, url.String(), extraParams, "file", "products_no_header.csv")
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	//client := &http.Client{}
	//_, err = client.Do(req)   // client.Do(req) just blocks in testing!
	handleGoodsCreateCc(w, req) //comment out handleGoodsList(w, r) & handleAccountError(...) in handler

	//t.Errorf("%s ~  ~ ", url.String(), w.Body)
	//t.Errorf("%s ~ %s ~ %s", url.String(), resp.Status, resp.Body)

}

/*func newFileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

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

	for k, v := range params {
		_ = writer.WriteField(k, v)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err

}*/

func newInstRequest(inst aetest.Instance, uri string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	go streamingUploadFile(params, paramName, path, w, file)
	req, err := inst.NewRequest("POST", uri, r)
	//req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	req.Header.Set("Content-Type", "multipart/form-data; boundary=doofus")
	return req, err
}

// Streams upload directly from file -> mime/multipart -> pipe -> http-request
func streamingUploadFile(params map[string]string, paramName, path string, w *io.PipeWriter, file *os.File) {
	defer file.Close()
	defer w.Close()
	writer := multipart.NewWriter(w)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
                return
		//return fmt.Errorf("Upload error: %v", err)
	}
	//sc <- writer.FormDataContentType()

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		//return fmt.Errorf("Upload error: %v", err)
		fmt.Println(err)
                return
	}
}

// Creates a new file upload http request with optional extra params
/*func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	go streamingUploadFile(params, paramName, path, w, file)
	return http.NewRequest("POST", uri, r)
}*/

func TestMain(t *testing.T) {
	/*path, _ := os.Getwd()
	path += "/test.pdf"
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest("http://localhost:8080", extraParams, "file", "test.pdf")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}*/
}
