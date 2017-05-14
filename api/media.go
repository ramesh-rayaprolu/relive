package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

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

func (api MediaAPI) transcodeMedia(absFileName, fullPath, fileName string) (err error) {
	var ffmpegPath string
	m3u8FileName := fmt.Sprintf("%s/%s.m3u8", fullPath, fileName)
	tsFileName := fmt.Sprintf("%s/%s%%d.ts", fullPath, fileName)

	ffmpegPath, err = exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Printf("Error looking up path for ffmpeg :%s\n", err.Error())
	}

	cmd := exec.Command(ffmpegPath,
		"-y", "-i", absFileName,
		"-codec", "copy", "-bsf", "h264_mp4toannexb",
		"-map", "0", "-f", "segment", "-segment_time", "10",
		"-segment_format", "mpegts", "-segment_list", m3u8FileName,
		"-segment_list_type", "m3u8", tsFileName)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start command, err: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Println("Exit Status: ", status.ExitStatus())
				if status.ExitStatus() != 0 {
					return fmt.Errorf("Exit Status: %d", status.ExitStatus())
				}
			}
		} else {
			return fmt.Errorf("error in  cmd.Wait: %v", err)
		}
	}
	return nil
}

// /api/media/play
func handleMediaPlayBack(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {

	fileToPlay := strings.TrimPrefix(args[0], "/api/media/play/")
	fileToPlay = fmt.Sprintf("/tmp/%s", fileToPlay)
	api.LogObj.PrintInfo("playing: %s", fileToPlay)

	http.ServeFile(w, r, fileToPlay)
	return nil
}

// /api/media/search
func handleMediaSearch(api MediaAPI, args []string, w http.ResponseWriter, r *http.Request) error {
	var id, pid uint64
	var fname string
	var err error
	var parsedURLSuffix *url.URL
	var result []dbmodel.MediaTypeEntry

	URLSuffix := args[0]
	parsedURLSuffix, err = url.Parse(URLSuffix)
	params := parsedURLSuffix.Query()

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

	if len(params["pid"]) > 0 {
		pid, err = strconv.ParseUint(params["pid"][0], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return fmt.Errorf("invalid parent id specified in request URL")
		}
	}

	if len(params["filename"]) > 0 {
		fname = params["filename"][0]
	}

	result, err = api.MediaDBI.SearchMediaTypeByID(id, pid, fname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return err
	}

	err = writeResponse(result, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
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

	fExt := filepath.Ext(header.Filename)
	fName := header.Filename[0 : len(header.Filename)-len(fExt)]

	outfilePath := fmt.Sprintf("/tmp/%s/%s", params["id"][0], fName)
	err = os.MkdirAll(outfilePath, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return err
	}

	outfileName := fmt.Sprintf("%s/%s", outfilePath, header.Filename)
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

	//file upload complete. transcode the file for smooth playback.
	err = api.transcodeMedia(outfileName, outfilePath, fName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return fmt.Errorf("Cannot transcode Media file: %v", err)
	}

	/*update DB only if all of the above succeeds */
	mediaURL := fmt.Sprintf("http://localhost:9999/api/media/play/%s/%s/%s.m3u8", params["id"][0], fName, fName)
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
	regex = "/api/media/search\\?([^/]+)$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaSearch,
		},
	)
	regex = "/api/media/play/([^/]+)/([^/]+)/([^/]+)$"
	media = append(media,
		mediaT{
			regex: regex,
			re:    regexp.MustCompile(regex),
			f:     handleMediaPlayBack,
		},
	)
}
