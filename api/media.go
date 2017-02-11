package api

import (
	"fmt"
	"net/http"
	"regexp"

	"../dbi"
)

// MediaAPI struct
type MediaAPI struct {
	MediaDBI dbi.MediaTypeTblDBI
}

//	/api/media/search
func handleMediaSearch(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// /api/media/store
func handleMediaStore(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
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
	regex = "/api/media/store$"
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
