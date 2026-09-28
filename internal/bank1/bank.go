// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

type withdrawal struct {
	amount int
	ok     chan bool
}

var withdrawals = make(chan withdrawal)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

// упр 9.1.
func Withdraw(amount int) bool {
	req := withdrawal{amount: amount, ok: make(chan bool)}
	withdrawals <- req
	return <-req.ok
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case req := <-withdrawals:
			successFlag := false
			if balance >= req.amount {
				balance -= req.amount
				successFlag = true
			}
			req.ok <- successFlag
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
