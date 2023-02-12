package main

import (
	"os"
	"testing"
)

func TestScrollDownText(t *testing.T) {
	tests := map[string]struct {
		file   string
		offset int
		height int

		want int
		ok   bool
	}{
		"Height=1|Offset=0": {
			file:   "testdata/288-lines.txt",
			offset: 0,
			height: 1,

			want: 1,
			ok:   true,
		},
		"Height=50|Offset=0": {
			file:   "testdata/288-lines.txt",
			offset: 0,
			height: 50,

			want: 25,
			ok:   true,
		},
		"Height=51|Offset=0": {
			file:   "testdata/288-lines.txt",
			offset: 0,
			height: 51,

			want: 26,
			ok:   true,
		},
		"Height=51|Offset=280": {
			file:   "testdata/288-lines.txt",
			offset: 280,
			height: 51,

			want: 0,
			ok:   false,
		},
		"Height=51|Offset=230": {
			file:   "testdata/288-lines.txt",
			offset: 230,
			height: 51,

			want: 237,
			ok:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			data := readFile(t, test.file)
			lines, err := splitByLines(data)
			if err != nil {
				t.Fatalf("splitting %s: %s", test.file, err)
			}

			got, ok := scrollDownText(lines, test.offset, test.height)

			if got != test.want || ok != test.ok {
				t.Errorf("got offset: %d, ok: %v, want: %d, %v", got, ok, test.want, test.ok)
			}
		})
	}
}

func TestScrollUpText(t *testing.T) {
	tests := map[string]struct {
		file   string
		offset int
		height int

		want int
		ok   bool
	}{
		"Height=1|Offset=0": {
			file:   "testdata/288-lines.txt",
			offset: 0,
			height: 1,

			want: 0,
			ok:   false,
		},
		"Height=50|Offset=0": {
			file:   "testdata/288-lines.txt",
			offset: 0,
			height: 50,

			want: 0,
			ok:   false,
		},
		"Height=51|Offset=288": {
			file:   "testdata/288-lines.txt",
			offset: 288,
			height: 51,

			want: 262,
			ok:   true,
		},
		"Height=50|Offset=200": {
			file:   "testdata/288-lines.txt",
			offset: 200,
			height: 50,

			want: 175,
			ok:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			data := readFile(t, test.file)
			lines, err := splitByLines(data)
			if err != nil {
				t.Fatalf("splitting %s: %s", test.file, err)
			}

			got, ok := scrollUpText(lines, test.offset, test.height)

			if got != test.want || ok != test.ok {
				t.Errorf("got offset: %d, ok: %v, want: %d, %v", got, ok, test.want, test.ok)
			}
		})
	}
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading %s: %s", filename, err)
	}

	return data
}
