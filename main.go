package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"main/download"
	"main/scrap"
	"os"
	"os/exec"
	"strconv"
)

//TODO/notes

// add checksum validation (checksums are on the website)
//
//	checksums still fail for some reason, probably they are moved in the csv file by some value for some reason
//
// implement automatic detection of new games (raise the "size" variable)
//
//	implement a .ini file that will store the latest "size"
//
// print speed of a download
//
// proxy https://stackoverflow.com/questions/33585587/creating-a-go-socks5-client
// workers https://twin.sh/articles/39/go-concurrency-goroutines-worker-pools-and-throttling-made-simple

const NumberOfWorkers = 1
const Csv_file_name = "stats.csv"
const Path_to_save = `Z:\gry`

func main() {
	if _, err := os.Stat(Path_to_save); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Path does not exist")
		return
	}
	download.Dow_from_file(NumberOfWorkers, Csv_file_name, Path_to_save)
	scrap.MakeCsvFile(Csv_file_name)
	//upload()
}

// TODO
//
// run the python file NOT in a sub-thread but just in the same terminal (fucky wacky)
//
// # OR
//
// get the credentials in go and just write a .ini file for the python script to use (little weird but should work just fine)
func upload() {
	exec.Command("ia_upload.py").Run()
}

// TODO
//
//	cache the output to .ini file
func calculateSize(csv_file string, exclude []string) {
	f, err := os.Open(csv_file)
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
			fmt.Println(err)
		}

		if rec[1] != "PlayStation 2" && rec[1] != "Xbox 360" && rec[1] != "PlayStation 3" && rec[1] != "Wii" && rec[1] != "WiiWare" && rec[1] != "Xbox" {
			b := rec[5]
			format := b[len(b)-2:]
			size, _ := strconv.ParseFloat(b[:len(b)-3], 32)

			if format == "KB" {
				mass += size
			} else if format == "MB" {
				mass += (size * 1024)
			} else if format == "GB" {
				mass += (size * (1024 * 1024))
			}
		}
	}
	fmt.Println("GB: ", mass/(1024*1024))
}

func setup_python() {
	//Creates a new python venv, if it already exists it does nothing
	//(python takes care of that)
	venv_cmd := exec.Command("python", "-m", "venv", ".")
	install_ia_cmd := exec.Command("pip", "install", "internetarchive")
	out, err := venv_cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	out, err = install_ia_cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
