package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
)

// https://github.com/gocolly/colly
type Info struct {
	platform string
	name     string
	download string
	year     string
	size     string
}

func main() {
	var infomainiac Info

	f, err := os.OpenFile("stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()

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

	//year
	c.OnXML("/html/body/div[4]/div[2]/div/div[3]/div[2]/div[1]/table/tbody/tr[3]/td[3]", func(e *colly.XMLElement) {
		infomainiac.year = e.Text
	})

	//size
	c.OnXML("//*[@id=\"download_size\"]", func(e *colly.XMLElement) {
		infomainiac.size = e.Text
	})

	for i := 1; i <= 5000; i++ {
		c.Visit("https://vimm.net/vault/" + strconv.Itoa(i))
		writer.Write([]string{strconv.Itoa(i), infomainiac.platform, infomainiac.name, "https:" + infomainiac.download + "?mediaId=" + strconv.Itoa(i), infomainiac.year, infomainiac.size})
		err = downloadFile("https:"+infomainiac.download+"?mediaId="+strconv.Itoa(i), infomainiac.download[10:11], infomainiac.platform)
		fmt.Println(err)
	}

}

func downloadFile(url string, ver string, platform string) (err error) {

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}

	req.Header.Set("Referer", "https://vimm.net/")
	req.Header.Set("Host", "download"+ver+".vimm.net")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//fmt.Println(resp.Header.Values("Content-Disposition"))
	filename := resp.Header.Values("Content-Disposition")[0]
	filename = filename[22 : len(filename)-1]
	//filename = strings.ReplaceAll(filename, " ", "_")
	defer resp.Body.Close()

	// Create the file
	path := "./" + platform + "/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	out, err := os.OpenFile("./"+platform+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
