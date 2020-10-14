package main

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/diff/fd"
)

func main() {
	// a "big" step
	run(options{
		Epsilon:      10e-5,
		LearningRate: 0.25,
		Method:       "fixedStep",
	})

	// a "small" step
	run(options{
		Epsilon:      10e-5,
		LearningRate: 0.01,
		Method:       "fixedStep",
	})

	// an "acceptable" step
	run(options{
		Epsilon:      10e-5,
		LearningRate: 0.24,
		Method:       "fixedStep",
	})

	// a "small" epsilon value
	run(options{
		Epsilon: 10e-5,
		Method:  "optimalStep",
	})

	// a "big" epsilon value
	run(options{
		Epsilon: 10e-2,
		Method:  "optimalStep",
	})
}

type options struct {
	Epsilon      float64
	LearningRate float64
	Method       string
}

func run(opts options) {
	// Use two variables instead of X = (x, y) to define the initial state.
	var x, y float64 = 1000, 1000
	var i int

	startTime := time.Now()

	for {
		var nextX, nextY float64

		if opts.Method == "fixedStep" {
			// xt+1 = xt - η * ∆xt
			nextX = x - opts.LearningRate*fd.Derivative(fx, x, nil)

			// yt+1 = yt - η * ∆yt
			nextY = y - opts.LearningRate*fd.Derivative(fy, y, nil)
		} else if opts.Method == "optimalStep" {
			// calculate the optimal step
			optimalLearningRate := (math.Pow(x, 2) + (math.Pow(7, 2) * math.Pow(y, 2))) / (math.Pow(x, 2) + (math.Pow(7, 3) * math.Pow(y, 2)))

			// xt+1 = xt - optimalLearningRate * ∆xt
			nextX = x - optimalLearningRate*fd.Derivative(fx, x, nil)

			// yt+1 = yt - optimalLearningRate * ∆yt
			nextY = y - optimalLearningRate*fd.Derivative(fy, y, nil)

			// store the optimal learning rate
			opts.LearningRate = optimalLearningRate
		}

		// Stop when (||f(xk+1) - f(xk)|| / ||f(xk)||) < Ɛ.
		if (math.Abs(f(nextX, nextY)-f(x, y)) / math.Abs(f(x, y))) < opts.Epsilon {
			break
		}

		// Store the new values into global variables.
		x, y = nextX, nextY
		i++
	}

	fmt.Printf("%+v, duration: %v, iterations: %d, x: %f, y: %f\n", opts, time.Since(startTime), i, x, y)
}

func f(x, y float64) float64 {
	return fx(x) + fy(y)
}

func fx(x float64) float64 {
	return 0.5 * math.Pow(x, 2)
}

func fy(y float64) float64 {
	return 3.5 * math.Pow(y, 2)
}
