package main

import (
	"log"
	"net"
	"time"

	"bufio"
	"fmt"

	"github.com/sony/gobreaker"
)

const addr = "127.0.0.1:5433"

func main() {

	// set up CB-parameters
	var st gobreaker.Settings
	st.Name = "Payment Server"
	st.MaxRequests = 2 // отправляем до двух запросов в состоянии half-open
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 3
	}
	st.Timeout = 10 * time.Second
	// create CB
	var cb *gobreaker.CircuitBreaker
	cb = gobreaker.NewCircuitBreaker(st)

	timeout := 1 * time.Second
	orders := []string{
		"123432-234",
		"123432-234",
		"123432-234",
		"853432-332",
		"853432-332",
		"853432-332",
		"254432-341",
		"254432-341",
		"254432-341",
		"254432-341",
		"853432-332",
		"853432-332",
		"853432-332",
		"254432-341",
		"123432-234",
		"123432-234",
		"853432-332",
		"123432-234",
	}

	for _, order := range orders {
		time.Sleep(2 * time.Second)
		println()
		// wrap client's business logic
		message, err := cb.Execute(func() (interface{}, error) {
			log.Print("trying to connect...")
			// connection request
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				err = fmt.Errorf("Connection error: %w", err)
				return nil, err
			}
			defer conn.Close()

			// pay-request
			message, err := pay(order, conn)
			if err != nil {
				err = fmt.Errorf("Payment error: %w", err)
				return nil, err
			}

			return message, nil
		})

		// logging error or message
		if err != nil {
			log.Printf("Service unavailable: %v\n\nCB state is: %s", err, cb.State().String())
		} else {
			log.Printf("Message from server: %s\nCB state is: %s", message, cb.State().String())
		}

	}
}

func pay(order string, conn net.Conn) (string, error) {
	fmt.Fprintf(conn, "please, process order %q\n", order)
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return message, nil
}
