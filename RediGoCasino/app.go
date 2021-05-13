package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// For docker:
// Without persistence: `docker run -p 6379:6379 --name redigocasino -d redis`
// With persistence: `docker run -p 6379:6379 --name redigocasino -d redis redis-server --appendonly yes`

// Check if name only consists of letters, numbers and space
func validName(name string) bool {
	const allowed = "abcdefghijklmnopqrstuvwxyz1234567890 "
	if len(name) <= 0 {
		return false
	}
	for _, char := range name {
		if !strings.Contains(allowed, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

func main() {
	// Initialization
	// Assumes redis is hosted locally at port 6379 (default port)
	connection, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().UnixNano())

	// Prompt name
	fmt.Println("Welcome to RediGo Casino!")
	name := ""
	for !validName(name) {
		fmt.Println("What is your name?")
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		name = strings.TrimSpace(name)
	}

	// Check if name exists
	// Idea: set authentication system (so that a password is needed before making bets)
	result, err := redis.Int(connection.Do("ZSCORE", "balance", name))
	if err == redis.ErrNil {
		// Name does not exist, initialize new account with a balance of 100
		result = 100
		fmt.Printf("Hello %s, this is probably your first time here.\n", name)
		_, err := connection.Do("ZADD", "balance", 100, name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Your balance is set to 100.")
	} else if err != nil {
		// Some other error occured
		log.Fatal(err)
	} else {
		// Name exists, print current balance
		if result == 0 {
			_, err := connection.Do("ZADD", "balance", 100, name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("It looks like you are broke, your balance will be resetted to 100.")
			result = 100
		}
		fmt.Printf("Hello %s, your current balance is %d.\n", name, result)
	}

	fmt.Println("The rules are simple, you have a 50% chance to double or lose your bet.")
	fmt.Println("How much do you want to bet?")

	// Prompt for bet, and validate input
	bet := -1
	for bet <= 0 || bet > result {
		scanint, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		scanint = strings.TrimSpace(scanint)
		convert, err := strconv.Atoi(scanint)
		if err != nil {
			continue
		}
		bet = convert
	}

	fmt.Println("You wagered", bet)
	// Casino logic, feel free to make it more complex
	if rand.Float64() <= 0.5 {
		_, err = connection.Do("ZINCRBY", "balance", -bet, name)
		if err != nil {
			log.Fatal(err)
		}
		// Lose
		fmt.Printf("You lost %d, your new balance is %d.\n", bet, result-bet)
	} else {
		// Win
		_, err = connection.Do("ZINCRBY", "balance", bet, name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You won %d, your new balance is %d.\n", bet, result+bet)
	}

	connection.Close()
	log.Println("Connection closed")
}
