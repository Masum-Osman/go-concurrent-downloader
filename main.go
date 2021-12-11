package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Download struct {
	Url          string
	TargetPath   string
	TotalSection int
}

func (d Download) Do() error {
	fmt.Println("Making Connections from DO Func...")
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	fmt.Printf("Got %v\n", resp.StatusCode)
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	fmt.Printf("Size of this file is: %v\n", size)

	var section = make([][2]int, d.TotalSection)
	eachSize := size / d.TotalSection
	fmt.Printf("Size of each section : %v\n", eachSize)
	fmt.Println("Section before shaped: ", section)

	for i := range section {
		if i == 0 {
			section[i][0] = 0
		} else {
			section[i][0] = section[i-1][1] + 1
		}

		if i < d.TotalSection-1 {
			section[i][1] = section[i][0] + eachSize
		} else {
			section[i][1] = size - 1
		}
	}

	fmt.Println("Section after shaped: ", section)

	var wg sync.WaitGroup
	for i, s := range section {
		wg.Add(1)
		i := i
		s := s
		go func() {
			defer wg.Done()
			err := d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()

	return d.mergeFiles(section)
}

func (d Download) downloadSection(i int, c [2]int) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c[0], c[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode))
	}
	fmt.Printf("Downloaded %v bytes for section %v %v\n", resp.Header.Get("Content-Length"), i, c)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (d Download) mergeFiles(section [][2]int) error {
	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range section {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return err
		}
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		err = os.Remove(tmpFileName)
		if err != nil {
			return err
		}
		fmt.Printf("%v bytes merged\n", n)
	}
	return nil
}

func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.Url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Silly Download Manager v001")
	return r, nil

}

func main() {
	startTime := time.Now()
	d := Download{
		Url:          "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4",
		TargetPath:   "final.mp4",
		TotalSection: 10,
	}

	err := d.Do()
	if err != nil {
		fmt.Println("An error occured while downloading the file: ", err)
	}

	fmt.Printf("Downloade completed in %v seconds\n", time.Now().Sub(startTime).Seconds())

}
