package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gocolly/colly/v2"
)

type Info struct {
	platform string
	name     string
	download string
	year     string
	size     string
	value    string
}

type proxyclient struct {
	client    *http.Client
	Transport http.Transport
}

const NumberOfWorkers = 1

func main() {
	reqch := make(chan *grab.Request)
	respch := make(chan *grab.Response, 10)

	f, err := os.Open("stats2.csv")
	if err != nil {
		fmt.Println(err)
	}
	csvReader := csv.NewReader(f)

	/*
		f2, err := os.OpenFile("downloads.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		writer := csv.NewWriter(f2)

		writer.Write([]string{"Id", "Platform", "Download", "Size"})
	*/

	proxyURL, err := url.Parse("http://localhost:9150")
	if err != nil {
		fmt.Println(err)
	}

	transport := http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	c := NewClient(transport)

	client := grab.Client{
		HTTPClient: c,
	}

	wg := sync.WaitGroup{}
	for i := 0; i < NumberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			client.DoChannel(reqch, respch)
			wg.Done()
		}()
	}

	start := time.Now()

	go func() {
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}

			b := rec[5]
			size, _ := strconv.ParseFloat(b[:len(b)-3], 32)
			if rec[1] != "PlayStation 2" && rec[1] != "Xbox 360" && rec[1] != "PlayStation 3" && rec[1] != "Wii" && rec[1] != "WiiWare" && rec[1] != "Xbox" && size > 0 {
				path := "./" + rec[1] + "/"
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir(path, os.ModePerm)
					if err != nil {
						log.Println(err)
					}
				}

				req, _ := grab.NewRequest(path, rec[3])
				//req.RateLimiter = rate.NewLimiter(rate.Limit(1048576), 1048576*2)
				req.HTTPRequest.Header.Set("Referer", "https://vimm.net/")
				req.HTTPRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0")

				reqch <- req

			}
		}
		close(reqch)
		wg.Wait()
		close(respch)
	}()

	for resp := range respch {
		// block until complete
		if err := resp.Err(); err != nil {
			panic(err)
		}

		fmt.Printf("Downloaded %s to %s\n", resp.Request.URL(), resp.Filename)
	}

	fmt.Printf("%s", time.Since(start))
	fmt.Println("")

	//updateCsvFile()
	//makeCsvFile()
	//downloadFile("https://download3.vimm.net/download/?mediaId=40342", "Game Boy Color")
}

// TODO
// https://twin.sh/articles/39/go-concurrency-goroutines-worker-pools-and-throttling-made-simple
/*
	for result := range jobResultChannel {
		jobResults = append(jobResults, result)
	}
*/
//fmt.Println(jobResults)

func (c *proxyclient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewClient return http client with a ratelimiter
func NewClient(trans http.Transport) *proxyclient {
	c := &proxyclient{
		client:    http.DefaultClient,
		Transport: trans,
	}
	return c
}

func calculateSize() {
	f, err := os.Open("stats2.csv")
	if err != nil {
		fmt.Println(err)
	}
	csvReader := csv.NewReader(f)

	var mass float64
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//if rec[1] == "Nintendo" || rec[2] == "Master System" || rec[2] == "Genesis" || rec[2] == "Super Nintendo" || rec[2] == "Nintendo 64" {
		if rec[1] != "PlayStation 2" && rec[1] != "Xbox 360" && rec[1] != "PlayStation 3" && rec[1] != "Wii" && rec[1] != "WiiWare" && rec[1] != "Xbox" {
			b := rec[5]
			format := b[len(b)-2:]
			size, _ := strconv.ParseFloat(b[:len(b)-3], 32)

			if format == "KB" {
				mass += size
			} else if format == "MB" {
				mass += (size * 1024)
			} else if format == "GB" {
				fmt.Println(b, size)
				mass += (size * (1024 * 1024))
				//fmt.Println(rec[1], rec[2], rec[5])
			}
		}
		//size, _ := strconv.Atoi(rec[5])
		//}
	}
	fmt.Println("GB: ", mass/(1024*1024))
}

func makeCsvFile() {

	//convert [email protected](cuz of cloudflare) to normal game name

	var infomainiac Info
	var status string

	f, err := os.OpenFile("stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(f)

	writer.Write([]string{"Id", "Platform", "Name", "Download", "Year", "Size"})

	c := colly.NewCollector(
		colly.AllowedDomains("vimm.net"),
	)

	//platform
	c.OnXML("/html/body/div[4]/div[2]/div/div[3]/h2/span[1]", func(e *colly.XMLElement) {
		infomainiac.platform = e.Text
	})

	//name
	c.OnXML("/html/body/div[4]/div[2]/div/div[3]/h2[1]/span[2]", func(e *colly.XMLElement) {
		infomainiac.name = e.Text
	})

	//download
	c.OnXML("//*[@id=\"download_form\"]", func(e *colly.XMLElement) {
		infomainiac.download = e.Attr("action")
	})

	c.OnXML("//*[@id=\"download_form\"]", func(e *colly.XMLElement) {
		infomainiac.value = e.ChildAttr("//input[1]", "value")
	})
	//year
	c.OnXML("/html/body/div[4]/div[2]/div/div[3]/div[2]/div[1]/table/tbody/tr[3]/td[3]", func(e *colly.XMLElement) {
		infomainiac.year = e.Text
	})

	//size
	c.OnXML("//*[@id=\"download_size\"]", func(e *colly.XMLElement) {
		infomainiac.size = e.Text
	})

	c.OnXML("/html/body/div[4]/div[2]/div/div[3]/p", func(e *colly.XMLElement) {
		status = e.Text
	})

	status = ""
	size := 88007
	for i := 1; i <= size; i++ {
		c.Visit("https://vimm.net/vault/" + strconv.Itoa(i))
		if status == "Error: Game not found." {
			status = ""
		} else {
			writer.Write([]string{strconv.Itoa(i), infomainiac.platform, infomainiac.name, "https:" + infomainiac.download + "?mediaId=" + infomainiac.value, infomainiac.year, infomainiac.size})
			writer.Flush()
			//fmt.Print(infomainiac.value, " || ")
			fmt.Println(strconv.Itoa(i) + "/" + strconv.Itoa(size))
		}
	}
}

func updateCsvFile() {
	f, err := os.Open("stats.csv")
	if err != nil {
		fmt.Println(err)
	}

	f2, err := os.Create("stats2.csv")
	if err != nil {
		fmt.Println(err)
	}

	writer := csv.NewWriter(f2)
	csvReader := csv.NewReader(f)

	writer.Write([]string{"Id", "Platform", "Name", "Download", "Year", "Size"})

	b := ""
	//i := 0
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if b != rec[2] {
			b = rec[2]
			writer.Write(rec)
			writer.Flush()
		}
		//b = rec[2]
	}
}

/*
func downloadFile(url string, platform string, workerid int) {
	// Create the file
	path := "./" + platform + "/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Get the data
	req, _ := grab.NewRequest("./"+platform+"/", url)
	client := grab.NewClient()

	resp := client.Do(req)
	//resp.Done (channel) check if open to check if the download is still going
	//add error handling
}
*/
