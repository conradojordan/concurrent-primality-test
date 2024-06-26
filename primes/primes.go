package primes

import (
	"context"
	"math"
	"sync"
)

type SearchInterval struct {
	start int
	end   int
}

func createSearchIntervals(possiblePrime int, numIntervals int) []SearchInterval {
	lowerLimit := 2
	upperLimit := int(math.Floor(math.Sqrt(float64(possiblePrime))))
	universe := upperLimit - lowerLimit + 1
	intervalSize, remainder := universe/numIntervals, universe%numIntervals
	intervals := make([]SearchInterval, 0, numIntervals)

	for j := 0; j < remainder; j++ {
		inter := SearchInterval{lowerLimit, lowerLimit + intervalSize}
		intervals = append(intervals, inter)
		lowerLimit += intervalSize + 1
	}
	for j := 0; j < (numIntervals - remainder); j++ {
		inter := SearchInterval{lowerLimit, lowerLimit + intervalSize - 1}
		intervals = append(intervals, inter)
		lowerLimit += intervalSize
	}
	return intervals
}

func isPrimeUnder1000(possiblePrime int) bool {
	if possiblePrime < 2 {
		return false
	}
	primesUnder1000 := [...]bool{2: true, 3: true, 5: true, 7: true, 11: true, 13: true, 17: true, 19: true, 23: true, 29: true, 31: true, 37: true, 41: true, 43: true, 47: true, 53: true, 59: true, 61: true, 67: true, 71: true,
		73: true, 79: true, 83: true, 89: true, 97: true, 101: true, 103: true, 107: true, 109: true, 113: true, 127: true, 131: true, 137: true, 139: true, 149: true, 151: true, 157: true, 163: true, 167: true, 173: true,
		179: true, 181: true, 191: true, 193: true, 197: true, 199: true, 211: true, 223: true, 227: true, 229: true, 233: true, 239: true, 241: true, 251: true, 257: true, 263: true, 269: true, 271: true, 277: true, 281: true,
		283: true, 293: true, 307: true, 311: true, 313: true, 317: true, 331: true, 337: true, 347: true, 349: true, 353: true, 359: true, 367: true, 373: true, 379: true, 383: true, 389: true, 397: true, 401: true, 409: true,
		419: true, 421: true, 431: true, 433: true, 439: true, 443: true, 449: true, 457: true, 461: true, 463: true, 467: true, 479: true, 487: true, 491: true, 499: true, 503: true, 509: true, 521: true, 523: true, 541: true,
		547: true, 557: true, 563: true, 569: true, 571: true, 577: true, 587: true, 593: true, 599: true, 601: true, 607: true, 613: true, 617: true, 619: true, 631: true, 641: true, 643: true, 647: true, 653: true, 659: true,
		661: true, 673: true, 677: true, 683: true, 691: true, 701: true, 709: true, 719: true, 727: true, 733: true, 739: true, 743: true, 751: true, 757: true, 761: true, 769: true, 773: true, 787: true, 797: true, 809: true,
		811: true, 821: true, 823: true, 827: true, 829: true, 839: true, 853: true, 857: true, 859: true, 863: true, 877: true, 881: true, 883: true, 887: true, 907: true, 911: true, 919: true, 929: true, 937: true, 941: true,
		947: true, 953: true, 967: true, 971: true, 977: true, 983: true, 991: true, 997: true, 1000: false,
	}
	return primesUnder1000[possiblePrime]
}

func hasDivisorInInterval(possiblePrime int, interval SearchInterval) bool {
	for i := interval.start; i <= interval.end; i++ {
		if possiblePrime%i == 0 {
			return true
		}
	}
	return false
}

func IsPrime(possiblePrime, numGoroutines int) bool {

	if possiblePrime <= 1000 {
		return isPrimeUnder1000(possiblePrime)
	}
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	ctx, cancel := context.WithCancel(context.Background())
	intervals := createSearchIntervals(possiblePrime, numGoroutines)
	foundDivisor := make(chan bool, numGoroutines)
	result := make(chan bool, 1)

	for i := 0; i < numGoroutines; i++ {
		go func(interval SearchInterval) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case foundDivisor <- hasDivisorInInterval(possiblePrime, interval):
				return
			}
		}(intervals[i])
	}

	go func() {
		for v := range foundDivisor {
			if v {
				cancel()
				result <- false
				return
			}
		}
		result <- true
	}()

	wg.Wait()
	close(foundDivisor)

	return <-result
}
