package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	// Статическое значение горутин
	const k = 5

	var q = flag.String("substring", "", "Substring to search in the input url")
	flag.Parse()

	result := make(chan int, k)
	end := make(chan struct{})

	//ограничение по горутинам
	gorCount := make(chan struct{}, k)

	// Запускаем Слушателя, который выведет результат по завершению
	go listener(result, end)

	wg := &sync.WaitGroup{}

	reader := bufio.NewScanner(os.Stdin)
	// Запускаем Пулл Воркеров, которые будут обрабатывать урлы по мере поступления
	for reader.Scan() {
		url := reader.Text()
		// Занимаем очередь
		gorCount <- struct{}{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			val, err := Worker(url, *q)
			if err != nil {
				log.Fatal(err)
			}

			result <- val

			//освобождаем очередь
			<-gorCount
		}()
	}
	wg.Wait()

	close(result)
	close(gorCount)

	<-end

}

func Worker(url string, q string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error on statusCode: %s", resp.Status)
		return -1, err
	}

	reader := resp.Body
	defer reader.Close()

	bd, _ := ioutil.ReadAll(reader)
	count := bytes.Count(bd, []byte(q))

	fmt.Printf("Count for %v: %d\n", url, count)

	return count, nil
}

func listener(result <-chan int, end chan struct{}) {

	count := 0
	for val := range result {
		count += val
	}
	fmt.Println("Total ", count)

	end <- struct{}{}
}
