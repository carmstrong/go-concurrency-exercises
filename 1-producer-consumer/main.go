//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func producer(stream Stream, c chan *Tweet, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			tweet, err := stream.Next()
			fmt.Println(tweet)
			if err == nil {
				c <- tweet
			} else {
				if err != ErrEOF {
					log.Fatal(err)
				}
				close(c)
			}
		}
	}()
}

func consumer(c chan *Tweet, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range c {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}
	}()
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	var c = make(chan *Tweet)
	var wg sync.WaitGroup

	// Producer
	go producer(stream, c, &wg)

	// Consumer
	go consumer(c, &wg)

	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
