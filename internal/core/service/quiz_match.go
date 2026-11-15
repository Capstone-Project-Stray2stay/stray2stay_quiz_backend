package service

import (
	"strings"
	"unicode"
)

// normalizeBreedText lowercases, strips punctuation (so "Shih-Tzu" and
// "shih tzu" compare equal), and collapses whitespace.
func normalizeBreedText(s string) string {
	var b strings.Builder
	lastWasSpace := false
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastWasSpace = false
		case unicode.IsSpace(r) || r == '-' || r == '\'' || r == '_':
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

// isCloseEnough reports whether guess matches breedName after normalization,
// tolerating a small number of typos proportional to word length so e.g.
// "siames" still matches "Siamese" but "persian" doesn't match "siamese".
func isCloseEnough(guess, breedName string) bool {
	g := normalizeBreedText(guess)
	b := normalizeBreedText(breedName)

	if g == "" {
		return false
	}
	if g == b {
		return true
	}

	distance := levenshtein(g, b)
	allowance := len(b) / 8
	if allowance < 1 {
		allowance = 1
	}

	return distance <= allowance
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			min := del
			if ins < min {
				min = ins
			}
			if sub < min {
				min = sub
			}
			curr[j] = min
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}
