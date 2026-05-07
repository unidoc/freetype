// Copyright 2026 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package truetype

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// FuzzParse exercises Parse and GlyphBuf.Load with arbitrary input.
// Seeds come from the bundled .ttf testdata; the fuzzer mutates them to
// surface format/bounds bugs (e.g. issue 598-style OOB reads on
// truncated or subsetted tables). A panic-free run for any input
// constitutes the success criterion - Parse and Load are allowed to
// return an error on malformed data, but never to panic.
func FuzzParse(f *testing.F) {
	seeds, err := filepath.Glob("../testdata/*.ttf")
	if err != nil {
		f.Fatalf("glob seed corpus: %v", err)
	}
	for _, p := range seeds {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		fnt, err := Parse(data)
		if err != nil {
			return
		}
		gb := &GlyphBuf{}
		scale := fixed.Int26_6(fnt.FUnitsPerEm())
		// Probe both in-range and clearly out-of-range indices so that
		// loca/glyf/hmtx bounds checks are exercised on every input.
		for _, i := range []Index{0, 1, 2, 32, 64, 1024, 65535} {
			_ = gb.Load(fnt, scale, i, font.HintingNone)
		}
	})
}
