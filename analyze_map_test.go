/**
 * Copyright (c) F5, Inc.
 *
 * This source code is licensed under the Apache License, Version 2.0 license found in the
 * LICENSE file in the root directory of this source tree.
 */

package crossplane

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestAnalyzeMapBody(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		mapDirective string
		parameter    *Directive
		term         string
		wantErr      *ParseError
	}{
		"valid map": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "default",
				Args:      []string{"0"},
				Line:      5,
				Block:     Directives{},
			},
			term: ";",
		},
		"valid map volatile parameter": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "volatile",
				Line:      5,
				Block:     Directives{},
			},
			term: ";",
		},
		"invalid map volatile parameter": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "volatile",
				Args:      []string{"1"},
				Line:      5,
				Block:     Directives{},
			},
			term:    ";",
			wantErr: &ParseError{What: "invalid number of parameters"},
		},
		"valid map hostnames parameter": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "hostnames",
				Line:      5,
				Block:     Directives{},
			},
			term: ";",
		},
		"invalid map hostnames parameter": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "hostnames",
				Args:      []string{"foo"},
				Line:      5,
				Block:     Directives{},
			},
			term:    ";",
			wantErr: &ParseError{What: "invalid number of parameters"},
		},
		"valid geo proxy_recursive parameter": {
			mapDirective: "geo",
			parameter: &Directive{
				Directive: "proxy_recursive",
				Line:      5,
				Block:     Directives{},
			},
			term: ";",
		},
		"invalid geo proxy_recursive parameter": {
			mapDirective: "geo",
			parameter: &Directive{
				Directive: "proxy_recursive",
				Args:      []string{"1"},
				Line:      5,
				Block:     Directives{},
			},
			term:    ";",
			wantErr: &ParseError{What: "invalid number of parameters"},
		},
		"valid geo ranges parameter": {
			mapDirective: "geo",
			parameter: &Directive{
				Directive: "ranges",
				Line:      5,
				Block:     Directives{},
			},
			term: ";",
		},
		"invalid geo ranges parameter": {
			mapDirective: "geo",
			parameter: &Directive{
				Directive: "ranges",
				Args:      []string{"0", "0", "0"},
				Line:      5,
				Block:     Directives{},
			},
			term:    ";",
			wantErr: &ParseError{What: "invalid number of parameters"},
		},
		"invalid number of parameters in map": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "default",
				Args:      []string{"0", "0"},
				Line:      5,
				Block:     Directives{},
			},
			term:    ";",
			wantErr: &ParseError{What: "invalid number of parameters"},
		},
		"missing semicolon": {
			mapDirective: "map",
			parameter: &Directive{
				Directive: "default",
				Args:      []string{"0", "0"},
				Line:      5,
				Block:     Directives{},
			},
			term:    "}",
			wantErr: &ParseError{What: `unexpected "}"`},
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := analyzeMapBody("nginx.conf", tc.parameter, tc.term, tc.mapDirective)
			if tc.wantErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)

			var perr *ParseError
			require.ErrorAs(t, err, &perr)
			require.Equal(t, tc.wantErr.What, perr.What)
		})
	}
}
