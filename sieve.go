package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println()
	flagTarget := flag.Int("n", 1000*1000, "number of primes to calculate")
	flagThreads := flag.Int("t", 2, "compute threads")
	flagWorkerSize := flag.Int("ws", 50000, "worker size")
	flag.Parse()

	target := *flagTarget
	threads := *flagThreads
	workSize := *flagWorkerSize
	if workSize%2 != 0 {
		workSize++
	}

	var j, prime, end, lastPrime int

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	initPrimes := []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
		31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
		73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
		127, 131, 137, 139, 149, 151, 157, 163, 167, 173,
		179, 181, 191, 193, 197, 199, 211, 223, 227, 229,
	}
	primes := make([]int, target+(workSize*threads))
	primeCounter := 0
	for _, initPrime := range initPrimes {
		primes[primeCounter] = initPrime
		primeCounter++
		lastPrime = initPrime
	}

	start := primes[primeCounter-1] + 2

	cand := start
	isPrime := true

	for ; primeCounter < 1000; cand += 2 {
		isPrime = true
		for _, prime = range primes {
			if prime == 0 || prime*prime > cand {
				break
			}
			for j = prime * ((cand / prime) - 1); j <= cand; j += prime {
				if cand == j {
					isPrime = false
					break
				}
			}
			if !isPrime {
				break
			}
		}
		if !isPrime {
			continue
		}
		primes[primeCounter] = cand
		primeCounter++
		lastPrime = cand
	}
	var wg2 sync.WaitGroup
	//var primeList *[]int
	for primeCounter <= target {
		var wg sync.WaitGroup
		wg.Add(threads)
		generatedPrimes := make([]*[]int, threads)

		for j := 0; j < threads; j++ {
			start = lastPrime + 2 + workSize*j
			end = lastPrime + workSize*(j+1)
			outPrimes := make([]int, workSize/5)
			generatedPrimes[j] = &outPrimes
			go sieve(&primes, &outPrimes, start, end, &wg)
		}
		wg.Wait()
		wg2.Wait()
		wg2.Add(1)
		addPrimes(&primes, &generatedPrimes, &primeCounter, &lastPrime, &wg2)
	}

	fmt.Println()
	fmt.Printf("Finished %s: %s\n", numberFormat(target), numberFormat(primes[target-1]))

}

func sieve(primeList, primesOut *[]int, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	isPrime := true
	cand := start
	var prime, i, j int
	for ; cand <= end; cand += 2 {
		isPrime = true
		for _, prime = range *primeList {
			if prime == 0 || prime*prime > cand {
				break
			}
			modifier := int(cand / prime)
			if modifier%2 == 0 {
				modifier--
			}
			for j = prime * modifier; j <= cand; j += 2 * prime {
				if cand == j {
					isPrime = false
					break
				}
			}
			if !isPrime {
				break
			}
		}
		if !isPrime {
			continue
		}
		(*primesOut)[i] = cand
		i++
	}
}

func addPrimes(outPrimes *[]int, inPrimes *[]*[]int, primeCounter, lastPrime *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, primeList := range *inPrimes {
		for j := 0; ; j++ {
			prime := (*primeList)[j]
			if prime == 0 {
				break
			}
			(*outPrimes)[*primeCounter] = prime
			*primeCounter++
		}
	}
	*lastPrime = (*outPrimes)[*primeCounter-1]
}

func numberFormat(number int) string {
	t := fmt.Sprintf("%d", number)
	i := 0
	o := ""
	for k := len(t) - 1; k >= 0; k-- {
		o += string(t[k])
		i++
		if k > 0 && i == 3 {
			o += ","
			i = 0
		}
	}
	t = ""
	for k := len(o) - 1; k >= 0; k-- {
		t += string(o[k])
	}
	return t
}
