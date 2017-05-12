package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/logger"
)

// MediaAPI struct
type MediaAPI struct {
	MediaDBI dbi.MediaTypeTblDBI
	LogObj   *logger.Logger
}

//	/api/media/search
func handleMediaSearch(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/media/store
func handleMediaStore(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	var formname, filename string
	URLSuffix := args[0]
	parsedURLSuffix, err := url.Parse(URLSuffix)
	params := parsedURLSuffix.Query()

	if len(params["formname"]) > 0 {
		formname = params["formname"][0]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required query parameters NOT specified in search request")
	}

	if len(params["filename"]) > 0 {
		filename = params["filename"][0]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required query parameters NOT specified in search request")
	}

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile(formname)

	if err != nil {
		fmt.Fprintln(w, err)
		return err
	}

	defer file.Close()

	outfileName := "/tmp/" + filename
	out, err := os.Create(outfileName)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return err
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
	return nil
}

// /api/media/update -
func handleMediaUpdate(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/media/delete -
func handleMediaDelete(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

type mediaT struct {
	regex string
	re    *regexp.Regexp
	f     func(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error
}

var media []mediaT

func (api MediaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, d := range media {
		if d.re.MatchString(r.URL.String()) {
			err := d.f(api, d.re.FindStringSubmatch(r.URL.String()), w, r)
			if err != nil {
				returnMessage := fmt.Sprintf("%v", err)
				w.Write([]byte(returnMessage))
			}
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("No match found.\n"))
}

func init() {
	var regex string
	regex = "/api/media/search$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaSearch,
		},
	)
	regex = "/api/media/store\\?([^/]+)$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaStore,
		},
	)
	regex = "/api/media/update$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaUpdate,
		},
	)
	regex = "/api/media/delete$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaDelete,
		},
	)
}
