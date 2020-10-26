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

func BenchmarkPrefixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match("urn:oss:*:200000456:myphotos/hangzhou/2015/*", "urn:oss:b:200000456:myphotos/hangzhou/2015/aaa")
	}
}

func BenchmarkPrefixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile("^urn:oss:.*:200000456:myphotos/hangzhou/2015/.*$")
	f := []byte("urn:oss:b:200000456:myphotos/hangzhou/2015/aaa")

	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
	}
}

func BenchmarkPrefixGlobMatch(b *testing.B) {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile("urn:oss:*:200000456:myphotos/hangzhou/2015/*")

	for i := 0; i < b.N; i++ {
		g.Match("urn:oss:b:200000456:myphotos/hangzhou/2015/aaa") // true
	}
}

func BenchmarkSuffixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match("*:oss:ListBuckets", "ucs:oss:ListBuckets")
	}
}
func BenchmarkSuffixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile("^.*:oss:ListBuckets")
	f := []byte("ucs:oss:ListBuckets")

	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
	}
}

func BenchmarkSuffixGlobMatch(b *testing.B) {
	var g glob.Glob

	// create simple glob
	g = glob.MustCompile("*:oss:ListBuckets")

	for i := 0; i < b.N; i++ {
		g.Match("ucs:oss:ListBuckets") // true
	}
}

func BenchmarkPrefixSuffixFilepathMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = filepath.Match(pattern_prefix_suffix, fixture_prefix_suffix_match)
	}
}
func BenchmarkPrefixSuffixRegexpMatch(b *testing.B) {
	m := regexp.MustCompile(regexp_prefix_suffix)
	f := []byte(fixture_prefix_suffix_match)
	// b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Match(f)
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
