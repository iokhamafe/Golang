/*
The classical Dining philosophers problem.

Implemented with forks (aka chopsticks) as mutexes.
*/

/*
Implement the dining philosopher’s problem with the following constraints/modifications.

There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
The host allows no more than 2 philosophers to eat concurrently.
Each philosopher is numbered, 1 through 5.

When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>”
on a line by itself, where <number> is the number of the philosopher.

When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>”
on a line by itself, where <number> is the number of the philosopher.
*/

package main

import (
	"fmt"
	"time"

	"sync"
)

var eatWgroup sync.WaitGroup

type fork struct{ sync.Mutex }

type philosopher struct {
	id                  int
	leftFork, rightFork *fork
}

// Goes from thinking to hungry to eating and done eating then starts over.
// Adapt the pause values to increased or decrease contentions
// around the forks.
func (p philosopher) eat() {
	defer eatWgroup.Done()
	for j := 0; j < 3; j++ {
		p.leftFork.Lock()
		p.rightFork.Lock()

		say("eating", p.id)
		time.Sleep(time.Second)

		p.rightFork.Unlock()
		p.leftFork.Unlock()

		say("finished eating", p.id)
		time.Sleep(time.Second)
	}

}

func say(action string, id int) {
	fmt.Printf("Philosopher #%d is %s\n", id+1, action)
}

func main() {
	// How many philosophers and forks

	count := 5

	// Create forks
	forks := make([]*fork, count)
	for i := 0; i < count; i++ {
		forks[i] = new(fork)
	}

	// Create philospoher, assign them 2 forks and send them to the dining table
	philosophers := make([]*philosopher, count)
	for i := 0; i < count; i++ {
		philosophers[i] = &philosopher{
			id: i, leftFork: forks[i], rightFork: forks[(i+1)%count]}
		eatWgroup.Add(1)
		go philosophers[i].eat()

	}
	eatWgroup.Wait()

}
