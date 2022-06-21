package main

import (
	"fmt"
	"time"
)

type Account struct {
	Name   string
	Amount float64
}

type BankOperation struct {
	Amount float64
	Done   chan struct{}
}

func transfer(amount float64, source, dest *Account) {
	if source.Amount < amount {
		fmt.Println("Error: No es posible realizar la transacción")
		return
	}

	time.Sleep(time.Second)

	source.Amount -= amount
	dest.Amount += amount

	fmt.Printf("Transacción exitosa %+v, %+v\n", source, dest)
}

func main() {
	signal := make(chan struct{})
	transaction := make(chan *BankOperation)

	client1 := Account{Name: "Juan", Amount: 500}
	client2 := Account{Name: "Pedro", Amount: 900}

	// Cajero
	go func() {
		for {
			request := <-transaction

			transfer(request.Amount, &client1, &client2)
			request.Done <- struct{}{}
		}
	}()

	operations := []float64{300.00, 300.00}
	for _, value := range operations {
		go func(amount float64) {
			requestTransaction := BankOperation{Amount: amount, Done: make(chan struct{})}
			transaction <- &requestTransaction

			signal <- <-requestTransaction.Done
		}(value)
	}

	<-signal
	<-signal
}
