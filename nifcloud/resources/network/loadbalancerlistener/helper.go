package loadbalancerlistener

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateLBImportString(importStr string) ([]string, error) {
	// example: example_8000_8000

	importParts := strings.Split(importStr, "_")
	errStr := "unexpected format of import string (%q), expected LBNAME_LBPORT_INSTANCEPORT: %s"
	if len(importParts) != 3 {
		return nil, fmt.Errorf(errStr, importStr, "invalid parts")
	}

	lbName := importParts[0]
	lbPort := importParts[1]
	instancePort := importParts[2]

	if lbName == "" {
		return nil, fmt.Errorf(errStr, importStr, "elb id must be required")
	}

	if _, err := strconv.Atoi(lbPort); err != nil {
		return nil, fmt.Errorf(errStr, importStr, "invalid lb port")
	}
	if _, err := strconv.Atoi(instancePort); err != nil {
		return nil, fmt.Errorf(errStr, importStr, "invalid instance port")
	}
	return importParts, nil
}

func populateLBFromImport(d *schema.ResourceData, importParts []string) error {
	lbName := importParts[0]
	lbPort := importParts[1]
	instancePort := importParts[2]

	p, err := strconv.Atoi(lbPort)
	if err != nil {
		return err
	}

	if err := d.Set("load_balancer_port", p); err != nil {
		return err
	}

	p, err = strconv.Atoi(instancePort)
	if err != nil {
		return err
	}

	if err := d.Set("instance_port", p); err != nil {
		return err
	}

	if err := d.Set("load_balancer_name", lbName); err != nil {
		return err
	}

	return nil
}

func getLBID(d *schema.ResourceData) string {
	return strings.Split(d.Id(), "_")[0]
}
