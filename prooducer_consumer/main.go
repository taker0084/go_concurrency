package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// NumberOfPizza is the max Count of pizza
const NumberOfPizzas = 10

var pizzaMade, pizzaFailed, total int

// Producer is a type for structs that holds two channels: one for pizza,
// with all information for a given pizza order including whether it was made
// successfully, and another to handle end og processing(when we quit the channel)
type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

// PizzaOrder is a type for structs that describes a given pizza order. It has that order
// number, a message indicating what happened to the order, and a boolean
// indicating if the order was successfully completed
type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

// Close is simply a method of closing the channel when we are done with it
// (i.e. something is pushed to the quit channel)
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

// makePizza attempts to make a pizza. We generate a random number from 1-12,
// and put in two cases where we can't make the pizza in time. Otherwise,
// we make the pizza without issue. To make things interesting, each pizza
// will take a different length of time to produce (some pizzas are harder than others).
func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		//rnd decides success or failure
		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzaFailed++
		} else {
			pizzaMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		//delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d ***!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d ***", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

// pizzeria is a goroutine that runs in the background and
// calls makePizza to try to make one order each time it iterates through
// the for loop. It executes until it receives something on the quit
// channel. The quit channel does not receive anything until the consumer
// sends it (when the number of orders is greater than or equal to the
// constant NumberOfPizzas).
func pizzeria(pizzaMaker *Producer) {
	//keep truck of which pizza we are making
	var i = 0

	//run forever or until we receive a quit notification
	//try to make pizzas
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			//we tried to make a pizza (we sent something to data channel -- a chan PizzaOrder)
			case pizzaMaker.data <- *currentPizza:

			//we want to quit, so send pizzaMaker.quit to the quitChan ( a chan error)
			case quitChan := <-pizzaMaker.quit:
				//close channel
				close(pizzaMaker.data)
				close(quitChan)
				close(pizzaMaker.quit)
				return
			}
		}
	}
}

func main() {
	//seed the random number generator
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))

	//print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	//create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	//run the producer in the background
	go pizzeria(pizzaJob)

	//create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	//print out the ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("we made %d pizzas, but failed to make %d, with %d attempts in total", pizzaMade, pizzaFailed, total)

	switch {
	case pizzaFailed > 9:
		color.Red("It was an awful day...")
	case pizzaFailed >= 6:
		color.Red("It was a pretty bad day")
	case pizzaFailed >= 4:
		color.Yellow("It was an okay day")
	case pizzaFailed >= 2:
		color.Yellow("It was a pretty good day")
	default:
		color.Green("It was a great day!")
	}
}
