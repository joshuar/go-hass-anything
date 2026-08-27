// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import (
	"strings"
	"testing"
)

func TestSanitiseClientID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "alphanumeric is preserved", id: "host123", want: "host123"},
		{name: "dots are dropped", id: "host.example.com", want: "hostexamplecom"},
		{name: "hyphens and underscores are dropped", id: "my-host_1", want: "myhost1"},
		{name: "non-ascii is dropped", id: "hößt", want: "ht"},
		{name: "empty stays empty", id: "", want: ""},
		{name: "only unusable characters becomes empty", id: "...---", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitiseClientID(tt.id); got != tt.want {
				t.Errorf("sanitiseClientID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestClientIDForHost(t *testing.T) {
	t.Run("is derived from the hostname", func(t *testing.T) {
		if got, want := clientIDForHost("somehost"), clientIDPrefix+"somehost"; got != want {
			t.Errorf("clientIDForHost() = %q, want %q", got, want)
		}
	})

	t.Run("is stable for the same hostname", func(t *testing.T) {
		if first, second := clientIDForHost("somehost"), clientIDForHost("somehost"); first != second {
			t.Errorf("clientIDForHost() is not stable: got %q then %q", first, second)
		}
	})

	t.Run("differs between hostnames", func(t *testing.T) {
		if a, b := clientIDForHost("hostone"), clientIDForHost("hosttwo"); a == b {
			t.Errorf("clientIDForHost() gave both hosts the same ID: %q", a)
		}
	})

	t.Run("hostnames differing only in punctuation still collide", func(t *testing.T) {
		// Documents a known limitation: punctuation is dropped, so "my-host"
		// and "myhost" share an ID. Hostnames that differ only in punctuation
		// are not expected in practice.
		if a, b := clientIDForHost("my-host"), clientIDForHost("myhost"); a != b {
			t.Errorf("expected %q and %q to match", a, b)
		}
	})

	t.Run("an unusable hostname yields a random ID", func(t *testing.T) {
		first := clientIDForHost("")
		if first == clientIDPrefix {
			t.Fatalf("clientIDForHost(%q) = %q, want an identifier after the prefix", "", first)
		}

		if second := clientIDForHost(""); first == second {
			t.Errorf("expected random IDs to differ, got %q twice", first)
		}
	})

	t.Run("a long hostname is hashed, not truncated", func(t *testing.T) {
		long := strings.Repeat("a", clientIDMaxHostLen) + "one"
		other := strings.Repeat("a", clientIDMaxHostLen) + "two"

		gotLong, gotOther := clientIDForHost(long), clientIDForHost(other)

		if gotLong == gotOther {
			t.Errorf("long hostnames sharing a prefix got the same ID: %q", gotLong)
		}

		if strings.Contains(gotLong, long) {
			t.Errorf("clientIDForHost() = %q, want the hostname hashed", gotLong)
		}

		if got, want := len(gotLong), len(clientIDPrefix)+16; got != want {
			t.Errorf("hashed client ID length = %d, want %d", got, want)
		}
	})
}
