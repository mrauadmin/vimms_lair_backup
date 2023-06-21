package scrap

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type Info struct {
	platform string
	name     string
	download string
	year     string
	size     string
	value    string
	md5      string
}

func MakeCsvFile(csv_file string) {
	a := 0
	if _, err := os.Stat(csv_file); err == nil {
		f, err := os.Open(csv_file)
		if err != nil {
			fmt.Println(err)
		}
		csvReader := csv.NewReader(f)
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			a, _ = strconv.Atoi(rec[0])
		}
		scrap_the_vimm(a+1, csv_file)
	} else if errors.Is(err, os.ErrNotExist) {
		//it needs to start from 1 as a page pf index 0 does not exist on vimm.net
		scrap_the_vimm(1, csv_file)
	} else {
		fmt.Println(err)
	}
}

func scrap_the_vimm(a int, csv_file string) {
	var infomainiac Info
	var infobuffer string

	f, err := os.OpenFile(csv_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(f)

	if a == 1 {
		writer.Write([]string{"Id", "Platform", "Name", "Download", "Year", "Size", "Md5"})
	}

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

	//download id
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

	//md5
	c.OnXML("//*[@id=\"data-md5\"]", func(e *colly.XMLElement) {
		infomainiac.md5 = e.Text
	})

	infobuffer = ""
	size := 88007
	for i := a; i <= size; i++ {
		c.Visit("https://vimm.net/vault/" + strconv.Itoa(i))
		if infobuffer != infomainiac.name {
			if infomainiac.name == "" {
				i++
			} else {
				infobuffer = infomainiac.name
				writer.Write([]string{strconv.Itoa(i), infomainiac.platform, infomainiac.name, "https:" + infomainiac.download + "?mediaId=" + infomainiac.value, infomainiac.year, infomainiac.size, infomainiac.md5})
				writer.Flush()

				infomainiac = flush_info_struct(infomainiac)

				fmt.Println(strconv.Itoa(i) + "/" + strconv.Itoa(size))
			}
		}
	}
}

func flush_info_struct(inf Info) Info {
	inf.platform = ""
	inf.name = ""
	inf.download = ""
	inf.year = ""
	inf.size = ""
	inf.value = ""
	inf.md5 = ""
	return inf
}
