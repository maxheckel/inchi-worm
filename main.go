package main

import (
	"fmt"
	"github.com/maxheckel/inchi-worm/model"
	"github.com/maxheckel/inchi-worm/utils"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const BaseURLFormat = "https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/InChIkey/%s/property/InChI/TXT"

func main() {

	argsWithoutProg := os.Args[1:]

	filePath := argsWithoutProg[0]
	lines, err := utils.ReadFileLines(filePath)
	if err != nil {
		panic(err)
	}

	res := []model.Inchi{}
	resChan := make(chan model.Inchi)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	fmt.Printf("Starting requests for %d keys", len(lines))

	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()
	for index, line := range lines {
		wg.Add(1)
		go getResultAsync(line, resChan, errChan, &wg)

		if index%5 == 0 {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
	}

	for result := range resChan {
		res = append(res, result)
	}
	for err := range errChan {
		fmt.Println(err)
	}

	utils.WriteOutput(res, "inchi_worm_output")
}

func getResultAsync(key string, resChan chan model.Inchi, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := getInchiValueForKey(key)
	if err != nil {
		errChan <- err
		return
	}
	resChan <- res
}

func getInchiValueForKey(key string) (model.Inchi, error) {
	url := fmt.Sprintf(BaseURLFormat, key)
	res, err := http.Get(url)
	if err != nil {
		return model.Inchi{}, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Inchi{}, err
	}
	bodyStr := string(bodyBytes)

	return model.Inchi{
		Key:   key,
		Value: strings.Split(bodyStr, "\n")[0],
	}, nil
}
