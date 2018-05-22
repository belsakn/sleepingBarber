package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Barber struct {
	name        string
	waitingRoom chan bool
	wake        chan bool
}

var waitingRoomReats int = 4
var wg sync.WaitGroup

func main() {

	newBarber := &Barber{name: "Joe", waitingRoom: make(chan bool, waitingRoomReats), wake: make(chan bool)}

	go startBarberShop(newBarber)
	//wait a bit for customers
	time.Sleep(time.Second * 2)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go customer(i, newBarber)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
	wg.Wait()
	fmt.Println("Bye...")
}

func startBarberShop(barber *Barber) {
	for {
		select {
		case <-barber.waitingRoom:
			fmt.Println("Start cutting...")
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
			fmt.Println("Stop cutting...")
		default:
			fmt.Println("let me sleep for some time")
			<-barber.wake
		}
	}
}

func customer(i int, barber *Barber) {
	select {
	case barber.waitingRoom <- true:
		fmt.Println("Client no.", i, "sent message")
		select {
		case barber.wake <- true:
			fmt.Println("Client no.", i, "waked him")
		default:
			fmt.Println("Client no.", i, "i will wait")
		}
		wg.Done()
		return

	default:
		fmt.Println("Client no.", i, "no message sent. I'm going home")
		time.Sleep(time.Millisecond * 100)
	}
	wg.Done()

}
