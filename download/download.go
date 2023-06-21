package download

import (
	"crypto/md5"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/cavaliergopher/grab/v3"
)

// TODO
//
// start_from
// start the downloads from this point
// it works but it isnt actually starting from correct point
// just compare strings of the last downloaded game and the strings of games in the csv file
func Dow_from_file(a_worker int, csv_file string, save_path string) {
	reqch := make(chan *grab.Request)
	respch := make(chan *grab.Response, 10)

	downfile, err := os.OpenFile("../downloaded.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	downloads := csv.NewWriter(downfile)

	logfile, err := os.OpenFile("../failure.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	logger := csv.NewWriter(logfile)

	f, err := os.Open(csv_file)
	if err != nil {
		fmt.Println(err)
	}
	csvReader := csv.NewReader(f)

	client := grab.NewClient()

	wg := sync.WaitGroup{}
	for i := 0; i < a_worker; i++ {
		wg.Add(1)
		go func() {
			client.DoChannel(reqch, respch)
			wg.Done()
		}()
	}

	start_from := scan_downloaded()
	go func() {
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}

			b := rec[5]
			size, _ := strconv.ParseFloat(b[:len(b)-3], 32)
			q, _ := strconv.Atoi(rec[0])
			if rec[1] != "PlayStation 2" && rec[1] != "Xbox 360" && rec[1] != "PlayStation 3" && rec[1] != "Wii" && rec[1] != "WiiWare" && rec[1] != "Xbox" && size > 0 && q > start_from {
				//Generates the folders
				path := save_path + "/" + rec[1] + "/"
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir(path, os.ModePerm)
					if err != nil {
						fmt.Println(err)
					}
				}
				req, _ := grab.NewRequest(path, rec[3])

				if rec[6] != "" {
					req.SetChecksum(md5.New(), []byte(rec[6]), true)
				}
				req.HTTPRequest.Header.Set("Referer", "https://vimm.net/")
				req.HTTPRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.71	0.05")
				reqch <- req
			}
		}
		close(reqch)
		wg.Wait()
		close(respch)
	}()

	b := scan_downloaded()
	for resp := range respch {
		err := resp.Err()
		if err != nil {
			logger.Write([]string{resp.Err().Error(), resp.Filename})
			logger.Flush()
		}
		downloads.Write([]string{strconv.Itoa(b), resp.Filename})
		downloads.Flush()
		b++
		fmt.Printf("Downloaded %s | Speed: %s \n", resp.Filename, strconv.FormatFloat(resp.Duration().Seconds(), 'f', 4, 64))
	}
}

func scan_downloaded() int {
	start_from := 0
	downfile_read, err := os.Open("../downloaded.csv")
	if err != nil {
		fmt.Println(err)
	}
	downloads_read := csv.NewReader(downfile_read)

	for {
		rec, err := downloads_read.Read()
		if err == io.EOF {
			break
		}
		start_from, _ = strconv.Atoi(rec[0])
	}

	return start_from
}
