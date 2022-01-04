package sock_pair_in_golang

import "math/rand"

func GenerateSocks(colors, patterns []string, numDuplicates int, onlySingles bool) Socks {
	socks := make(Socks, 0)
	if numDuplicates < 1 {
		return socks
	}

	for _, pattern := range patterns {
		for _, color := range colors {
			for i := 0; i < numDuplicates; i++ {
				if onlySingles {
					socks = append(
						socks,
						Sock{color, pattern, true},
					)
				} else {
					socks = append(
						socks,
						Sock{color, pattern, true},
						Sock{color, pattern, false},
					)
				}
			}
		}
	}

	return socks
}

func ShuffleSocks(socks Socks) Socks {
	rand.Shuffle(len(socks), func(i, j int) {
		socks[i], socks[j] = socks[j], socks[i]
	})

	return socks
}
