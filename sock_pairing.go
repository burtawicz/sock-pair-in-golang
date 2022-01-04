package sock_pair_in_golang

import (
	"errors"
	"math/rand"
	"sort"
)

// removeSockFromBasket removes the Sock at the specified index and returns the remaining list.
// If the index is out of bounds an error is returned.
func removeSockFromBasket(s Socks, idx int) (Socks, error) {
	if idx >= len(s) || idx < 0 {
		return s, errors.New("invalid index")
	}

	if len(s) == 1 && idx == 0 {
		return make([]Sock, 0), nil
	}

	if idx+1 <= len(s)-1 {
		return append(s[:idx], s[idx+1:]...), nil
	}

	return s[:idx], nil
}

// orderSockPair returns the pair of Sock ordered as left, right.
func orderSockPair(s1, s2 Sock) (Sock, Sock) {
	if s1.IsLeft {
		return s1, s2
	}
	return s2, s1
}

type SockPairingStrategy interface {
	pairSocks(freshSocks Socks) (SockPairs, Socks)
}

// RandomPairingStrategy is the process of grabbing the first Sock in the basket, then
// drawing a second random Sock from the basket for comparison.
type RandomPairingStrategy struct{}

func (s RandomPairingStrategy) pairSocks(freshSocks Socks) (SockPairs, Socks) {
	pairedSocks := make(SockPairs, 0)
	orphanedSocks := make(Socks, 0)

	// check for empty basket
	if len(freshSocks) == 0 {
		return pairedSocks, orphanedSocks
	}

	// check for single item basket
	if len(freshSocks) == 1 {
		return pairedSocks, append(orphanedSocks, freshSocks[0])
	}

	sockToPair := freshSocks[0]
	freshSocks = freshSocks[1:]
	reassignAndCheckForOrphans := false

	comparisonCount := 0
	for len(freshSocks) > 0 {
		// generate a random index
		randomIdx := rand.Intn(len(freshSocks))
		foundSock := freshSocks[randomIdx]

		if sockToPair.IsMatchingPair(foundSock) {
			// add our socks to the pairedSocks
			leftSock, rightSock := orderSockPair(sockToPair, foundSock)
			pairedSocks = append(pairedSocks, Socks{leftSock, rightSock})

			// remove the matched sock
			if res, err := removeSockFromBasket(freshSocks, randomIdx); err == nil {
				freshSocks = res
			}

			reassignAndCheckForOrphans = true
		} else if comparisonCount > len(freshSocks)*len(freshSocks) {
			orphanedSocks = append(orphanedSocks, sockToPair)
			reassignAndCheckForOrphans = true
		}

		comparisonCount++

		if reassignAndCheckForOrphans {
			if len(freshSocks) > 1 {
				// assign a new sockToPair to compare against
				sockToPair = freshSocks[0]
				// remove the new sockToPair from the unpairedSocks
				if res, err := removeSockFromBasket(freshSocks, 0); err == nil {
					freshSocks = res
				}
				// reset the index, so we can compare our sockToPair against all the unpairedSocks
				reassignAndCheckForOrphans = false
				// reset comparisonCount
				comparisonCount = 0
			} else {
				// last sock is an orphan
				orphanedSocks = append(orphanedSocks, freshSocks...)
				break
			}
		}
	}

	return pairedSocks, orphanedSocks
}

// SequentialPairingStrategy is the process of grabbing the first Sock in the basket, then
// comparing it to each subsequent Sock from the basket for comparison.
type SequentialPairingStrategy struct{}

func (s SequentialPairingStrategy) pairSocks(freshSocks Socks) (SockPairs, Socks) {
	pairedSocks := make(SockPairs, 0)
	orphanedSocks := make(Socks, 0)

	// check for empty basket
	if len(freshSocks) == 0 {
		return pairedSocks, orphanedSocks
	}

	// check for single item basket
	if len(freshSocks) == 1 {
		return pairedSocks, append(orphanedSocks, freshSocks[0])
	}

	sockToPair := freshSocks[0]
	freshSocks = freshSocks[1:]
	reassignAndCheckForOrphans := false

	i := 0
	for len(freshSocks) > 0 {
		if sockToPair.IsMatchingPair(freshSocks[i]) {
			// add our socks to the pairedSocks
			leftSock, rightSock := orderSockPair(sockToPair, freshSocks[i])
			pairedSocks = append(pairedSocks, Socks{leftSock, rightSock})

			// remove the matched sock
			if res, err := removeSockFromBasket(freshSocks, i); err == nil {
				freshSocks = res
			}

			reassignAndCheckForOrphans = true
		} else {
			i++

			// check if we've encountered an orphaned sock that is not at the bottom of the basket
			if i >= len(freshSocks) {
				// add sockToPair to orphanedSocks
				orphanedSocks = append(orphanedSocks, sockToPair)
				reassignAndCheckForOrphans = true
			}
		}

		if reassignAndCheckForOrphans {
			if len(freshSocks) > 1 {
				// assign a new sockToPair to compare against
				sockToPair = freshSocks[0]
				// remove the new sockToPair from the unpairedSocks
				if res, err := removeSockFromBasket(freshSocks, 0); err == nil {
					freshSocks = res
				}
				// reset the index, so we can compare our sockToPair against all the unpairedSocks
				i = 0
				reassignAndCheckForOrphans = false
			} else {
				// last sock is an orphan
				orphanedSocks = append(orphanedSocks, freshSocks...)
				break
			}
		}
	}

	return pairedSocks, orphanedSocks
}

// SortFirstPairingStrategy is the process of sorting all the socks in the basket by color, then
// comparing each nth and nth+1 Sock in the basket.
type SortFirstPairingStrategy struct{}

func (s SortFirstPairingStrategy) pairSocks(freshSocks Socks) (SockPairs, Socks) {
	pairedSocks := make(SockPairs, 0)
	orphanedSocks := make(Socks, 0)

	sort.Sort(freshSocks)

	i := 0
	for len(freshSocks) > 0 {
		if i+1 < len(freshSocks) {
			if freshSocks[i].IsMatchingPair(freshSocks[i+1]) {
				leftSock, rightSock := orderSockPair(freshSocks[i], freshSocks[i+1])
				pairedSocks = append(pairedSocks, Socks{leftSock, rightSock})
				i++
			} else {
				orphanedSocks = append(orphanedSocks, freshSocks[i])
			}
			// update socks
			freshSocks = freshSocks[i+1:]
			i = 0
		} else {
			orphanedSocks = append(orphanedSocks, freshSocks[i])
			freshSocks = freshSocks[i+1:]
		}
	}

	return pairedSocks, orphanedSocks
}

// SurfacePairingStrategy is the process of placing each Sock on a surface and checking for matches
// each time a new Sock is pulled from the basket.
type SurfacePairingStrategy struct{}

func (s SurfacePairingStrategy) pairSocks(freshSocks Socks) (SockPairs, Socks) {
	pairedSocks := make(SockPairs, 0)
	orphanedSocks := make(Socks, 0)
	surface := make(map[Sock]Socks)

	for _, sock := range freshSocks {
		matchingSock := Sock{sock.Color, sock.Pattern, !sock.IsLeft}
		// check if matching sock already exists on the surface
		if surface[matchingSock] != nil && len(surface[matchingSock]) > 0 {
			leftSock, rightSock := orderSockPair(sock, matchingSock)
			pairedSocks = append(pairedSocks, Socks{leftSock, rightSock})
			// remove the matching sock from the surface
			if res, err := removeSockFromBasket(surface[matchingSock], 0); err == nil {
				surface[matchingSock] = res
			}
		} else {
			if surface[sock] == nil {
				surface[sock] = make(Socks, 0)
			}
			surface[sock] = append(surface[sock], sock)
		}
	}

	// collect remaining orphaned socks
	for _, socks := range surface {
		orphanedSocks = append(orphanedSocks, socks...)
	}

	return pairedSocks, orphanedSocks
}
