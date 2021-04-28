package urlutil

import (
	"testing"
)

type sameHostTest struct {
	url1, url2 string
	same       bool
}

type joinTest struct {
	base, rel, res string
}

var (
	sameHostTests = []sameHostTest{
		{"http://www.gondolaweb.com", "", true},
		{"http://www.gondolaweb.com", "/", true},
		{"http://www.gondolaweb.com", "/foo", true},
		{"http://www.gondolaweb.com", "../anotherhost.com", true},
		{"http://www.gondolaweb.com", "//gondola", false},
		{"//gondola", "//gondola", true},
		{"//www.gondolaweb.com", "//gondola.com", false},
		{"http://www.gondolaweb.com", "https://twitter.com", false},
		{"http://www.gondolaweb.com", "https://www.gondolaweb.com", true},
	}

	joinTests = []joinTest{
		{"http://gondola", "foo", "http://gondola/foo"},
		{"http://gondola/bar", "foo", "http://gondola/foo"},
		{"http://gondola/bar/", "foo", "http://gondola/bar/foo"},
		{"http://gondola/bar/", "/foo/", "http://gondola/foo/"},
		{"//gondola/bar/", "/foo/", "//gondola/foo/"},
	}
)

func TestSameHost(t *testing.T) {
	for _, v := range sameHostTests {
		same1 := SameHost(v.url1, v.url2)
		same2 := SameHost(v.url2, v.url1)
		if same1 == same2 {
			if same1 != v.same {
				t.Errorf("expecting SameHost(%q, %q) = %v, got %v instead", v.url1, v.url2, v.same, same1)
			}
		} else {
			t.Errorf("non conmmutative SameHost(%q, %q) != SameHost(%q, %q)", v.url1, v.url2, v.url2, v.url1)
		}
	}
}

func TestJoin(t *testing.T) {
	for _, v := range joinTests {
		res, err := Join(v.base, v.rel)
		if err != nil {
			t.Errorf("error joining %q and %q: %s", v.base, v.rel, err)
			continue
		}
		if res != v.res {
			t.Errorf("expecting Join(%q, %q) = %q, got %q instead", v.base, v.rel, v.res, res)
		}
	}
}
