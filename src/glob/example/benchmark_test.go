package example

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gobwas/glob"
)

const (
	pattern_prefix                 = "abc*"
	regexp_prefix                  = `^abc.*$`
	pattern_suffix                 = "*def"
	regexp_suffix                  = `^.*def$`
	pattern_prefix_suffix          = "ab*ef"
	regexp_prefix_suffix           = `^ab.*ef$`
	fixture_prefix_suffix_match    = "abcdef"
	fixture_prefix_suffix_mismatch = "af"
)



func BenchmarkPrefixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile("^aaa:bbb:.*:cccccc:myphotos/hangzhou/2015/.*$")
	f := []byte("aaa:bbb:b:cccccc:myphotos/hangzhou/2015/aaa")

	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
	}
}

func BenchmarkPrefixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match("aaa:bbb:*:cccccc:myphotos/hangzhou/2015/*", "aaa:bbb:b:cccccc:myphotos/hangzhou/2015/aaa")
	}
}

func BenchmarkPrefixGlobMatch(b *testing.B) {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile("aaa:bbb:*:cccccc:myphotos/hangzhou/2015/*")

	for i := 0; i < b.N; i++ {
		g.Match("aaa:bbb:b:cccccc:myphotos/hangzhou/2015/aaa") // true
	}
}


func BenchmarkSuffixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile("^.*:aaa:abcabcabc")
	f := []byte("123:aaa:abcabcabc")

	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
	}
}

func BenchmarkSuffixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match("*:aaa:abcabcabc", "123:aaa:abcabcabc")
	}
}

func BenchmarkSuffixGlobMatch(b *testing.B) {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile("*:aaa:abcabcabc")

	for i := 0; i < b.N; i++ {
		g.Match("123:aaa:abcabcabc") // true
	}
}


func BenchmarkPrefixSuffixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile(regexp_prefix_suffix)
	f := []byte(fixture_prefix_suffix_match)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
	}
}

func BenchmarkPrefixSuffixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match(pattern_prefix_suffix, fixture_prefix_suffix_match)
	}
}

func BenchmarkPrefixSuffixGlobMatch(b *testing.B) {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile(pattern_prefix_suffix)

	for i := 0; i < b.N; i++ {
		g.Match(fixture_prefix_suffix_match) // true
	}
}

// go test -bench=. benchmark_test.go
