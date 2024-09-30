package toolkit

// execute tests: `go test .` or `o test -v .`
import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Error("wrong length random string returned")
	}

}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/jpeg", "image/png"}, renameFile: true, errorExpected: false},
	{name: "not allowed", allowedTypes: []string{"image/jpeg"}, renameFile: true, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		// set up a pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			/// create the form data field 'file'
			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()

		// read from the pipe which recieves data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			// clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
			// if err != nil {
			// 	t.Error("error removing file", err)
			// }
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", e.name)
		}

		wg.Wait()
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	// set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		/// create the form data field 'file'
		part, err := writer.CreateFormFile("file", "./testdata/img.png")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/img.png")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Error("error decoding image", err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	// read from the pipe which recieves data
	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools
	// testTools.AllowedFileTypes = uploadTests[0].allowedTypes

	uploadedFiles, err := testTools.UploadOneFile(request, "./testdata/uploads/", true)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}

	// clean up
	_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))
}

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools Tools

	err := testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	err = os.Remove("./testdata/myDir")

	if err != nil {
		t.Error(err)
	}
}

var slugTests = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "Now it the time 123", expected: "now-it-the-time-123", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "Now is the time for all GOOD men! + fish & such &~^123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
	{name: "japanese string", s: "こんにちは世界", expected: "", errorExpected: true},
	{name: "japanese string and roman characters", s: "Hello world こんにちは世界", expected: "hello-world", errorExpected: false},
}

func TestTools_Slugify(t *testing.T) {
	var testTool Tools

	for _, e := range slugTests {
		slug, err := testTool.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when not expected: %s", e.name, err.Error())
		}

		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug returned; expected %s but got %s", e.name, e.expected, slug)
		}
	}
}

func TestTools_DownloadStaticFile(t *testing.T) {
	resRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "./testdata", nil)

	var testTool Tools

	testTool.DownloadStaticFiles(resRecorder, req, "./testdata", "gopher.png", "downloaded.png")

	res := resRecorder.Result()
	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "226510" {
		t.Error("wrong content length of", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"downloaded.png\"" {
		t.Error("wrong content disposition", res.Header["Content-Disposition"][0])
	}

	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}

var jsonTests = []struct {
	name          string
	json          string
	errorExpected bool
	maxSize       int
	allowUnknown  bool
}{
	{name: "good jason", json: `{"foo":"bar"}`, errorExpected: false, maxSize: 1024, allowUnknown: false},
	{name: "badly formatted json", json: `{"foo":}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "incorrect type", json: `{"foo": 1}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "two json files", json: `{"foo": "1"}{"alpha": "beta"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "empty body", json: ``, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "syntax error in json", json: `{"foo": 1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "unknown field in json", json: `{"fooo": "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: false},
	{name: "unknown field in json", json: `{"fooo": "1"}`, errorExpected: false, maxSize: 1024, allowUnknown: true},
	{name: "missing field name", json: `{foo: "1"}`, errorExpected: true, maxSize: 1024, allowUnknown: true},
	{name: "file too large", json: `{"foo": "bar"}`, errorExpected: true, maxSize: 5, allowUnknown: true},
	{name: "not json", json: `Hello, world!`, errorExpected: true, maxSize: 1024, allowUnknown: true},
}

func TestTools_ReadJSON(t *testing.T) {
	var testTool Tools

	for _, e := range jsonTests {
		// set the max file size
		testTool.MaxJSONSize = e.maxSize

		// allow/disallow unknown fields
		testTool.AllowUnkownFields = e.allowUnknown

		// declare a variable to read the decoded json into
		var decodedJSON struct {
			Foo string `json:"foo"`
		}

		// create a request with the body
		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(e.json)))

		if err != nil {
			t.Log("Error:", err)
		}

		// create a recorder
		resRecorder := httptest.NewRecorder()

		err = testTool.ReadJSON(resRecorder, req, &decodedJSON)

		if e.errorExpected && err == nil {
			t.Errorf("%s: error expected, but none received", e.name)
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: error not expected, but one received: %s", e.name, err.Error())
		}

		req.Body.Close()
	}
}
