// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fetch

import (
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/pkgsite/internal/derrors"
	"golang.org/x/pkgsite/internal/proxy"
	"golang.org/x/pkgsite/internal/version"
)

func TestDirectoryModuleGetterEmpty(t *testing.T) {
	g, err := NewDirectoryModuleGetter("", "testdata/has_go_mod")
	if err != nil {
		t.Fatal(err)
	}
	if want := "example.com/testmod"; g.modulePath != want {
		t.Errorf("got %q, want %q", g.modulePath, want)
	}

	_, err = NewDirectoryModuleGetter("", "testdata/no_go_mod")
	if !errors.Is(err, derrors.BadModule) {
		t.Errorf("got %v, want BadModule", err)
	}
}

func TestEscapedPath(t *testing.T) {
	for _, test := range []struct {
		path, version, suffix string
		want                  string
	}{
		{
			"m.com", "v1", "info",
			"dir/cache/download/m.com/@v/v1.info",
		},
		{
			"github.com/aBc", "v2.3.4", "zip",
			"dir/cache/download/github.com/a!bc/@v/v2.3.4.zip",
		},
	} {
		g, err := NewFSProxyModuleGetter("dir")
		if err != nil {
			t.Fatal(err)
		}
		got, err := g.escapedPath(test.path, test.version, test.suffix)
		if err != nil {
			t.Fatal(err)
		}
		want, err := filepath.Abs(test.want)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s, %s, %s: got %q, want %q", test.path, test.version, test.suffix, got, want)
		}
	}
}

func TestFSProxyGetter(t *testing.T) {
	ctx := context.Background()
	const (
		modulePath = "github.com/jackc/pgio"
		vers       = "v1.0.0"
		goMod      = "module github.com/jackc/pgio\n\ngo 1.12\n"
	)
	ts, err := time.Parse(time.RFC3339, "2019-03-30T17:04:38Z")
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewFSProxyModuleGetter("testdata/modcache")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("info", func(t *testing.T) {
		got, err := g.Info(ctx, modulePath, vers)
		if err != nil {
			t.Fatal(err)
		}
		want := &proxy.VersionInfo{Version: vers, Time: ts}
		if !cmp.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}

		// Asking for latest should give the same version.
		got, err = g.Info(ctx, modulePath, version.Latest)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}

		if _, err := g.Info(ctx, "nozip.com", vers); !errors.Is(err, derrors.NotFound) {
			t.Errorf("got %v, want NotFound", err)
		}
	})
	t.Run("mod", func(t *testing.T) {
		got, err := g.Mod(ctx, modulePath, vers)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte(goMod)
		if !cmp.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if _, err := g.Mod(ctx, "nozip.com", vers); !errors.Is(err, derrors.NotFound) {
			t.Errorf("got %v, want NotFound", err)
		}
	})
	t.Run("contentdir", func(t *testing.T) {
		fsys, err := g.ContentDir(ctx, modulePath, vers)
		if err != nil {
			t.Fatal(err)
		}
		// Just check that the go.mod file is there and has the right contents.
		f, err := fsys.Open("go.mod")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		got, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte(goMod)
		if !cmp.Equal(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

		if _, err := g.ContentDir(ctx, "nozip.com", vers); !errors.Is(err, derrors.NotFound) {
			t.Errorf("got %v, want NotFound", err)
		}
	})
}
