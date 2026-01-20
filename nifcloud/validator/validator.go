package validator

import (
	"net"
	"unicode/utf8"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringRuneCountBetween returns a SchemaValidateFunc which tests if the provided value
// is of type string and has rune count between min and max (inclusive)
func StringRuneCountBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(v interface{}, k cty.Path) diag.Diagnostics {
		value, ok := v.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}

		if utf8.RuneCountInString(value) < min || utf8.RuneCountInString(value) > max {
			return diag.Errorf("expected length of %s to be in the range (%d - %d), got %s", k, min, max, value)
		}
		return nil
	}
}

// CIDRNetworkAddress returns a diag.Diagnostics which tests if the provided value
// is of type string and the string value is a valid CIDR that represents a network address
func CIDRNetworkAddress(v interface{}, k cty.Path) diag.Diagnostics {
	cidr, ok := v.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}

	ip1, ipnet1, err := net.ParseCIDR(cidr)
	if err != nil {
		return diag.Errorf("%q is not a valid CIDR block; did you mean %q?", cidr, ipnet1)
	}

	ip2, ipnet2, err := net.ParseCIDR(ipnet1.String())
	if err != nil {
		return diag.Errorf("%q is not a valid CIDR block; did you mean %q?", cidr, ipnet1)
	}

	if ip2.String() != ip1.String() || ipnet2.String() != ipnet1.String() {
		return diag.Errorf("%q is not a valid CIDR block; did you mean %q?", cidr, ipnet1)
	}

	return nil
}

// IPAddress returns a diag.Diagnostics which tests if the provided value
// is of type string and the string value is a valid IPAddress
func IPAddress(v interface{}, k cty.Path) diag.Diagnostics {
	s, ok := v.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return diag.Errorf("%q is not a valid IPAddress", s)
	}
	return nil
}

// Any TODO to use terraform-plugin-sdk after flow issue merged
// https://github.com/hashicorp/terraform-plugin-sdk/issues/534
func Any(validators ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(v interface{}, k cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, validator := range validators {
			d := validator(v, k)
			if len(d) == 0 {
				return nil
			}
			diags = append(diags, d...)
		}
		return diags
	}
}
