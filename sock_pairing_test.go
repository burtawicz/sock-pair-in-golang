package sock_pair_in_golang

import (
	"reflect"
	"sort"
	"testing"
)

type testCase struct {
	name        string
	freshSocks  Socks
	wantPairs   SockPairs
	wantOrphans Socks
}

func getTestCases() []testCase {
	return []testCase{
		{
			"3 matching pairs",
			Socks{
				Sock{"red", "plain", true},
				Sock{"green", "plain", true},
				Sock{"red", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", false},
			},
			SockPairs{
				Socks{Sock{"blue", "plain", true}, Sock{"blue", "plain", false}},
				Socks{Sock{"green", "plain", true}, Sock{"green", "plain", false}},
				Socks{Sock{"red", "plain", true}, Sock{"red", "plain", false}},
			},
			make([]Sock, 0),
		},
		{
			"3 matching pairs, 1 orphaned sock",
			Socks{
				Sock{"red", "plain", true},
				Sock{"green", "plain", true},
				Sock{"red", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", false},
				Sock{"pink", "plain", true},
			},
			SockPairs{
				Socks{Sock{"blue", "plain", true}, Sock{"blue", "plain", false}},
				Socks{Sock{"green", "plain", true}, Sock{"green", "plain", false}},
				Socks{Sock{"red", "plain", true}, Sock{"red", "plain", false}},
			},
			Socks{{"pink", "plain", true}},
		},
		{
			"all orphaned freshSocks",
			Socks{
				Sock{"red", "plain", true},
				Sock{"green", "plain", true},
				Sock{"red", "plain", true},
				Sock{"blue", "plain", true},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
				Sock{"pink", "plain", true},
			},
			make(SockPairs, 0),
			Socks{
				Sock{"blue", "plain", true},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
				Sock{"green", "plain", true},
				Sock{"pink", "plain", true},
				Sock{"red", "plain", true},
				Sock{"red", "plain", true},
			},
		},
		{
			"no socks",
			make(Socks, 0),
			make(SockPairs, 0),
			make(Socks, 0),
		},
		{
			"single sock",
			Socks{Sock{"pink", "plain", true}},
			make(SockPairs, 0),
			Socks{Sock{"pink", "plain", true}},
		},
	}
}

func Test_removeSockFromBasket(t *testing.T) {
	tests := []struct {
		name       string
		idx        int
		freshSocks Socks
		want       Socks
		wantErr    bool
	}{
		{
			"slice from 0:",
			0,
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			Socks{
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			false,
		},
		{
			"slice from 0:2,3:",
			2,
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			false,
		},
		{
			"slice from 0:5",
			5,
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
			},
			false,
		},
		{
			"slice from 0:4, 5:",
			4,
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"green", "plain", true},
			},
			false,
		},
		{
			"invalid index",
			15,
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			Socks{
				Sock{"red", "plain", true},
				Sock{"red", "plain", false},
				Sock{"green", "plain", false},
				Sock{"blue", "plain", false},
				Sock{"blue", "plain", true},
				Sock{"green", "plain", true},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removeSockFromBasket(tt.freshSocks, tt.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeSockFromBasket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeSockFromBasket() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderSockPair(t *testing.T) {
	tests := []struct {
		name      string
		s1        Sock
		s2        Sock
		wantLeft  Sock
		wantRight Sock
	}{
		{
			"left, right",
			Sock{"red", "plain", true},
			Sock{"red", "plain", false},
			Sock{"red", "plain", true},
			Sock{"red", "plain", false},
		},
		{
			"right, left",
			Sock{"red", "plain", false},
			Sock{"red", "plain", true},
			Sock{"red", "plain", true},
			Sock{"red", "plain", false},
		},
		{
			"left, left",
			Sock{"red", "plain", true},
			Sock{"red", "plain", true},
			Sock{"red", "plain", true},
			Sock{"red", "plain", true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeft, gotRight := orderSockPair(tt.s1, tt.s2)
			if !reflect.DeepEqual(gotLeft, tt.wantLeft) {
				t.Errorf("orderSockPair() gotLeft = %v, want %v", gotLeft, tt.wantLeft)
			}
			if !reflect.DeepEqual(gotRight, tt.wantRight) {
				t.Errorf("orderSockPair() gotRight = %v, want %v", gotRight, tt.wantRight)
			}
		})
	}
}

func TestRandomPairingStrategy_pairSocks(t *testing.T) {
	for _, tt := range getTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			strategy := RandomPairingStrategy{}
			gotPairs, gotOrphans := strategy.pairSocks(tt.freshSocks)
			// sort the returned values, so they can be compared to our shared test values
			sort.Sort(gotPairs)
			sort.Sort(gotOrphans)

			if !reflect.DeepEqual(gotPairs, tt.wantPairs) {
				t.Errorf("RandomPairingStrategy.pairSocks() pairs = %v, want %v", gotPairs, tt.wantPairs)
			}

			if !reflect.DeepEqual(gotOrphans, tt.wantOrphans) {
				t.Errorf("RandomPairingStrategy.pairSocks() orphans = %v, want %v", gotOrphans, tt.wantOrphans)
			}
		})
	}
}

func BenchmarkRandomPairingStrategy_pairSocks_noOrphans(b *testing.B) {
	b.Skip("Skip in short mode")
	strategy := RandomPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		5,
		false,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkRandomPairingStrategy_pairSocks_singleOrphan(b *testing.B) {
	b.Skip("Skip in short mode")
	strategy := RandomPairingStrategy{}
	testSocks := ShuffleSocks(append(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		5,
		false,
	), Sock{"pink", "plain", true}))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkRandomPairingStrategy_pairSocks_allOrphans(b *testing.B) {
	b.Skip("Skip in short mode")
	strategy := RandomPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		5,
		true,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func TestSequentialPairingStrategy_pairSocks(t *testing.T) {
	for _, tt := range getTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			strategy := SequentialPairingStrategy{}
			gotPairs, gotOrphans := strategy.pairSocks(tt.freshSocks)
			// sort the returned values, so they can be compared to our shared test values
			sort.Sort(gotPairs)
			sort.Sort(gotOrphans)

			if !reflect.DeepEqual(gotPairs, tt.wantPairs) {
				t.Errorf("SequentialPairingStrategy.pairSocks() pairs = %v, want %v", gotPairs, tt.wantPairs)
			}

			if !reflect.DeepEqual(gotOrphans, tt.wantOrphans) {
				t.Errorf("SequentialPairingStrategy.pairSocks() orphans = %v, want %v", gotOrphans, tt.wantOrphans)
			}
		})
	}
}

func BenchmarkSequentialPairingStrategy_pairSocks_noOrphans(b *testing.B) {
	strategy := SequentialPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSequentialPairingStrategy_pairSocks_singleOrphan(b *testing.B) {
	strategy := SequentialPairingStrategy{}
	testSocks := ShuffleSocks(append(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	), Sock{"pink", "plain", true}))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSequentialPairingStrategy_pairSocks_allOrphans(b *testing.B) {
	strategy := SequentialPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		true,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func TestSortFirstPairingStrategy_pairSocks(t *testing.T) {
	for _, tt := range getTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			strategy := SortFirstPairingStrategy{}
			gotPairs, gotOrphans := strategy.pairSocks(tt.freshSocks)
			// sort the returned values, so they can be compared to our shared test values
			sort.Sort(gotPairs)
			sort.Sort(gotOrphans)

			if !reflect.DeepEqual(gotPairs, tt.wantPairs) {
				t.Errorf("SortFirstPairingStrategy.pairSocks() pairs = %v, want %v", gotPairs, tt.wantPairs)
			}

			if !reflect.DeepEqual(gotOrphans, tt.wantOrphans) {
				t.Errorf("SortFirstPairingStrategy.pairSocks() orphans = %v, want %v", gotOrphans, tt.wantOrphans)
			}
		})
	}
}

func BenchmarkSortFirstPairingStrategy_pairSocks_noOrphans(b *testing.B) {
	strategy := SortFirstPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSortFirstPairingStrategy_pairSocks_singleOrphan(b *testing.B) {
	strategy := SortFirstPairingStrategy{}
	testSocks := ShuffleSocks(append(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	), Sock{"pink", "plain", true}))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSortFirstPairingStrategy_pairSocks_allOrphans(b *testing.B) {
	strategy := SortFirstPairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		true,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func TestSurfacePairingStrategy_pairSocks(t *testing.T) {
	for _, tt := range getTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			strategy := SurfacePairingStrategy{}
			gotPairs, gotOrphans := strategy.pairSocks(tt.freshSocks)
			// sort the returned values, so they can be compared to our shared test values
			sort.Sort(gotPairs)
			sort.Sort(gotOrphans)

			if !reflect.DeepEqual(gotPairs, tt.wantPairs) {
				t.Errorf("SurfacePairingStrategy.pairSocks() pairs = %v, want %v", gotPairs, tt.wantPairs)
			}

			if !reflect.DeepEqual(gotOrphans, tt.wantOrphans) {
				t.Errorf("SurfacePairingStrategy.pairSocks() orphans = %v, want %v", gotOrphans, tt.wantOrphans)
			}
		})
	}
}

func BenchmarkSurfacePairingStrategy_pairSocks_noOrphans(b *testing.B) {
	strategy := SurfacePairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSurfacePairingStrategy_pairSocks_singleOrphan(b *testing.B) {
	strategy := SurfacePairingStrategy{}
	testSocks := ShuffleSocks(append(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		false,
	), Sock{"pink", "plain", true}))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}

func BenchmarkSurfacePairingStrategy_pairSocks_allOrphans(b *testing.B) {
	strategy := SurfacePairingStrategy{}
	testSocks := ShuffleSocks(GenerateSocks(
		[]string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"},
		[]string{"plain", "checkered", "herringbone", "plaid", "striped"},
		150,
		true,
	))
	for i := 0; i < b.N; i++ {
		strategy.pairSocks(testSocks)
	}
}
