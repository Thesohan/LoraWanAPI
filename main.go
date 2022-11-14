package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// cleanup wait for all running goroutines to send the completion signal. It also
// close the channel
func cleanup(wg *sync.WaitGroup, channel chan bool) {
	wg.Wait()
	close(channel)
	fmt.Println("cleanup: all running goroutines are completed.")
}

func main() {
	result := make(map[int]string)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) //To capture the singles from OS

	channel := make(chan bool, ThrottleLimit) // To allow concurrent execution of at max ThrottleLimit goroutine
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating a goroutine to capture the signals from OS
	go func() {
		<-sigs
		// listen for OS signals
		fmt.Println(" received C-c - shutting down")
		// tell other goroutines to stop
		fmt.Println("telling other goroutines to stop")
		cancel() // To tell other goroutins to stop
		cleanup(&wg, channel)
	}()

	lorawanService := LorawanService{}

	for count := 0; count < LorawanBatchSize; count++ {
		channel <- true // if channel is full it will wait for till it allow the new value
		wg.Add(1)       // Increment number of running go routine
		go registerNewDevice(count, channel, ctx, lorawanService, &wg, result)
	}
	cleanup(&wg, channel)
	fmt.Println(result)

}

// registerNewDevice makes sure that a new device EUI is registered in the Lorawan.
func registerNewDevice(count int, channel chan bool, ctx context.Context, lorawanService LorawanService, wg *sync.WaitGroup, result map[int]string) {
	defer func() {
		wg.Done()
		<-channel
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(count, "goroutine ended")
			return
		default:
			deveui, err := GenerateHexString(HexStringLenght)
			// Retry if device is not registered.
			if err == nil && lorawanService.
				RegisterDevice(
					LorawanAPIRequestBody{
						Deveui: deveui,
					},
					count,
				) {
				fmt.Println(count, ": ", deveui)
				result[count] = deveui
				return
			}
		}
	}
}
