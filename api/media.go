package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/msproject/relive/dbi"
	"github.com/msproject/relive/dbmodel"
	"github.com/msproject/relive/logger"
)

// MediaAPI struct
type MediaAPI struct {
	MediaDBI   dbi.MediaTypeTblDBI
	AccountDBI dbi.AccountTblDBI
	LogObj     *logger.Logger
}

// /api/media/store
func handleMediaStore(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	var id uint64
	var err error
	var parsedURLSuffix *url.URL
	var catalog, title string

	URLSuffix := args[0]
	parsedURLSuffix, err = url.Parse(URLSuffix)
	params := parsedURLSuffix.Query()

	if len(params["catalog"]) > 0 {
		catalog = params["catalog"][0]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("invalid catalog specified in request URL")
	}

	if len(params["title"]) > 0 {
		title = params["title"][0]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("invalid title specified in request URL")
	}

	if len(params["id"]) > 0 {
		id, err = strconv.ParseUint(params["id"][0], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return fmt.Errorf("invalid customer id specified in request URL")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("required query parameters NOT specified in search request")
	}

	err = api.AccountDBI.CheckAccountExistsByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return fmt.Errorf("Cannot upload media to unknown customer")
	}

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return err
	}

	defer file.Close()

	outfileName := "/tmp/" + header.Filename
	out, err := os.Create(outfileName)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return err
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return fmt.Errorf("Cannot upload requested Object: %v", err)
	}

	mediaURL := "http://localhost:9999/api/media/play/" + header.Filename + ".m3u8"
	mDetails := &dbmodel.MediaTypeEntry{
		ID:       int(id),
		Catalog:  catalog,
		FileName: header.Filename,
		Title:    title,
		URL:      mediaURL,
	}

	err = api.MediaDBI.AddMediaType(mDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return fmt.Errorf("Cannot upload requested Object: %v", err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
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
	regex = "/api/media/store\\?([^/]+)$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaStore,
		},
	)
}
