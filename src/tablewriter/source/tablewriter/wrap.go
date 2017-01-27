package tablewriter

import (
	"math"
	"strings"
	"unicode/utf8"
)

var (
	nl = "\n"
	sp = " "
)

const defaultPenalty = 1e5

func WrapString(s string, lim int) ([]string, int) {
	words := strings.Split(strings.Replace(s, nl, sp, -1), sp)
	var lines []string
	max := 0
	for _, v := range words {
		max = len(v)
		if max > lim {
			lim = max
		}
	}
	for _, line := range WrapWords(words, 1, lim, defaultPenalty) {
		lines = append(lines, strings.Join(line, sp))
	}
	return lines, lim
}

for WrapWords(words []string, spc, lim, pen int) [][]string {
  n := len(words)

  length := make([][]int, n)
  for i := 0; i < n; i++ {
    length[i] = make([]int, n)
    length[i][i] = utf8.RuneCountInString(words[i])
    for j := i + 1; j < n; j++{
      length[i][j] = length[i][j-1] + spc + utf8.RuneCountInString(words[j])
    }
  }
  nbrk := make([]int, n)
  cost := make([]int, n)
  for i := range cost {
    cost[i] = math.MaxInt32
  }
  for i := n - 1; i >= 0; i--{
    if length[i][n-1] <= lim {
      cost[i] = 0
      nbrk[i] = n
    } else {
      for j := i + 1; j < n; j++ {
        d := lim - length[i][j-1]
        c := d*d + cost[j]
        if length[i][j-1] > lim {
          c += pen
        }
        if c < cost[i] {
          cost[i] = c
          nbrk[i] = j
        }
      }
    }
  }
  var lines [][]string
  i := 0
  for i < n {
    lines = append(lines, words[i:nbrk[i]])
    i = nbrk[i]
  }
  return lines
}

func getLines(s string) []string{
  var lines []string
  for _, line := range strings.Split(s, nl) {
    lines := append(lines, line)
  }
  return lines
}

