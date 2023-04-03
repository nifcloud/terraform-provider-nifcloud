package record

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateDnsRecordImportString(importStr string) ([]string, error) {
	// example: setIdentifier_example.test_@
	// example: setIdentifier_example.test_sub
	// example: setIdentifier_example.test_sub.example.test

	importParts := strings.Split(importStr, "_")
	errStr := "unexpected format of import string (%q), expected SETIDENTIFIER_ZONEID_NAME: %s"
	if len(importParts) < 3 {
		return nil, fmt.Errorf(errStr, importStr, "invalid parts")
	}

	setIdentifier := importParts[0]
	zoneID := importParts[1]
	name := importParts[2]

	if name == "" {
		return nil, fmt.Errorf(errStr, importStr, "name must be required")
	}

	if zoneID == "" {
		return nil, fmt.Errorf(errStr, importStr, "zone_id must be required")
	}

	if setIdentifier == "" {
		return nil, fmt.Errorf(errStr, importStr, "set_identifier must be required")
	}

	return importParts, nil
}

func populateDnsRecordFromImport(d *schema.ResourceData, importParts []string) error {
	setIdentifier := importParts[0]
	zoneID := importParts[1]
	name := strings.Join(importParts[2:], "_")

	if err := d.Set("set_identifier", setIdentifier); err != nil {
		return err
	}

	d.SetId(setIdentifier)

	if err := d.Set("zone_id", zoneID); err != nil {
		return err
	}

	if err := d.Set("name", name); err != nil {
		return err
	}

	return nil
}
