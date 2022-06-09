package pathmatch_test

import (
	"regexp"
	"testing"

	"github.com/abichinger/pathmatch"
)

func BenchmarkCompile(b *testing.B) {
	benchmarks := []struct {
		regex string
		path  string
	}{
		{
			`foo/([^/]+)/bar/(.*)`,
			"foo/:id/bar/*",
		},
	}

	for _, benchmark := range benchmarks {
		b.Run("regexp", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = regexp.Compile(benchmark.regex)
			}
		})
		b.Run("pathmatch", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = pathmatch.Compile(benchmark.path)
			}
		})
	}
}

func BenchmarkMatchString(b *testing.B) {
	benchmarks := []struct {
		regex string
		path  string
		str   string
	}{
		{
			`foo/([^/]+)/bar/(.*)`,
			"foo/:id/bar/*",
			"foo/1/bar/2/baz/3",
		},
	}

	for _, benchmark := range benchmarks {
		b.Run("regexp", func(b *testing.B) {
			r, err := regexp.Compile(benchmark.regex)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = r.MatchString(benchmark.str)
				//_ = r.FindStringSubmatch(benchmark.str)
			}
		})
		b.Run("pathmatch", func(b *testing.B) {
			p, err := pathmatch.Compile(benchmark.path)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.Match(benchmark.str)
				//_ = p.FindSubmatch(benchmark.str)
			}
		})
	}
}

func BenchmarkFindSubmatch(b *testing.B) {
	benchmarks := []struct {
		regex string
		path  string
		str   string
	}{
		{
			`foo/([^/]+)/bar/(.*)`,
			"foo/:id/bar/*",
			"foo/1/bar/2/baz/3",
		},
	}

	for _, benchmark := range benchmarks {
		b.Run("regexp", func(b *testing.B) {
			r, err := regexp.Compile(benchmark.regex)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = r.FindStringSubmatch(benchmark.str)
			}
		})
		b.Run("pathmatch", func(b *testing.B) {
			p, err := pathmatch.Compile(benchmark.path)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.FindSubmatch(benchmark.str)
			}
		})
	}
}
