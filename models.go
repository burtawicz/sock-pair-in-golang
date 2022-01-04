package sock_pair_in_golang

type Sock struct {
	Color   string
	Pattern string
	IsLeft  bool
}

func (s *Sock) IsMatchingPair(s2 Sock) bool {
	return s.Color == s2.Color &&
		s.Pattern == s2.Pattern &&
		s.IsLeft != s2.IsLeft
}

type Socks []Sock

func (s Socks) Len() int {
	return len(s)
}

func (s Socks) Less(i, j int) bool {
	s1, s2 := s[i], s[j]
	if s1.Color != s2.Color {
		return s1.Color < s2.Color
	}

	if s1.Pattern != s2.Pattern {
		return s1.Pattern < s2.Pattern
	}

	return s1.IsLeft
}

func (s Socks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SockPairs []Socks

func (s SockPairs) Len() int {
	return len(s)
}

func (s SockPairs) Less(i, j int) bool {
	return Socks{s[i][0], s[j][0]}.Less(0, 1)
}

func (s SockPairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
