package pathmatch_test

import (
	"regexp"
	"testing"

	"github.com/abichinger/pathmatch"
)

func BenchmarkCompile(b *testing.B) {
	benchmarks := []struct {
		name  string
		regex string
		path  string
	}{
		{
			"normal",
			`foo/([^/]+)/bar/(.*)`,
			"foo/:id/bar/*",
		},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name+"-regex", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = regexp.Compile(benchmark.regex)
			}
		})
		b.Run(benchmark.name+"-path", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = pathmatch.Compile(benchmark.path)
			}
		})
	}
}

func BenchmarkMatchString(b *testing.B) {
	benchmarks := []struct {
		name  string
		regex string
		path  string
		str   string
	}{
		{
			"normal",
			`foo/([^/]+)/bar/(.*)`,
			"foo/:id/bar/*",
			"foo/1/bar/2/baz/3",
		},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name+"-regex", func(b *testing.B) {
			r, err := regexp.Compile(benchmark.regex)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = r.FindStringSubmatch(benchmark.str)
			}
		})
		b.Run(benchmark.name+"-path", func(b *testing.B) {
			p, err := pathmatch.Compile(benchmark.path)
			if err != nil {
				b.Errorf(err.Error())
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.MatchString(benchmark.str)
			}
		})
	}
}
