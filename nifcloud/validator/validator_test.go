package validator

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func TestStringRuneCountBetween(t *testing.T) {
	cases := map[string]struct {
		Value         interface{}
		ExpectedDiags diag.Diagnostics
	}{
		"NotString": {
			Value:         7,
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"Empty": {
			Value:         "",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"NotWithinMin": {
			Value:         "あ",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"NotWithinMax": {
			Value:         "ああああ",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"BothWithin": {
			Value:         "あああ",
			ExpectedDiags: nil,
		},
	}

	fn := StringRuneCountBetween(2, 3)

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			diags := fn(tc.Value, cty.Path{})

			if len(diags) != len(tc.ExpectedDiags) {
				t.Fatalf("%s: wrong number of diags, expected %d, got %d", tn, len(tc.ExpectedDiags), len(diags))
			}

			for j := range diags {
				if diags[j].Severity != tc.ExpectedDiags[j].Severity {
					t.Fatalf("%s: expected severity %v, got %v", tn, tc.ExpectedDiags[j].Severity, diags[j].Severity)
				}
			}
		})
	}
}

func TestIPAddress(t *testing.T) {
	cases := map[string]struct {
		Value         interface{}
		ExpectedDiags diag.Diagnostics
	}{
		"NotString": {
			Value:         7,
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"Empty": {
			Value:         "",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"NotWithIPAddress": {
			Value:         "a",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"IpAddress": {
			Value:         "192.0.2.0",
			ExpectedDiags: nil,
		},
		"Ipv6Address": {
			Value:         "2001:db8::",
			ExpectedDiags: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			diags := IPAddress(tc.Value, cty.Path{})

			if len(diags) != len(tc.ExpectedDiags) {
				t.Fatalf("%s: wrong number of diags, expected %d, got %d", tn, len(tc.ExpectedDiags), len(diags))
			}

			for j := range diags {
				if diags[j].Severity != tc.ExpectedDiags[j].Severity {
					t.Fatalf("%s: expected severity %v, got %v", tn, tc.ExpectedDiags[j].Severity, diags[j].Severity)
				}
			}
		})
	}
}

func TestCIDRNetworkAddress(t *testing.T) {
	cases := map[string]struct {
		Value         interface{}
		ExpectedDiags diag.Diagnostics
	}{
		"NotString": {
			Value:         7,
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"Empty": {
			Value:         "",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"NotWithCIDR": {
			Value:         "a",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"IpAddress": {
			Value:         "192.0.2.0",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"InvalidCIDR": {
			Value:         "192.0.2.1/24",
			ExpectedDiags: diag.Diagnostics{{Severity: diag.Error}},
		},
		"ValidCIDR": {
			Value:         "192.0.2.0/24",
			ExpectedDiags: nil,
		},
		"ValidCIDRIPV6": {
			Value:         "2001:db8::/122",
			ExpectedDiags: nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			diags := CIDRNetworkAddress(tc.Value, cty.Path{})

			if len(diags) != len(tc.ExpectedDiags) {
				t.Fatalf("%s: wrong number of diags, expected %d, got %d", tn, len(tc.ExpectedDiags), len(diags))
			}

			for j := range diags {
				if diags[j].Severity != tc.ExpectedDiags[j].Severity {
					t.Fatalf("%s: expected severity %v, got %v", tn, tc.ExpectedDiags[j].Severity, diags[j].Severity)
				}
			}
		})
	}
}
