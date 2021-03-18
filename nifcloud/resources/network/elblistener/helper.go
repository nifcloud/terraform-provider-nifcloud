package elblistener

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/internal/mutexkv"
)

var mutexKV = mutexkv.NewMutexKV()

func validateELBImportString(importStr string) ([]string, error) {
	// example: example_TCP_8000_8000

	importParts := strings.Split(importStr, "_")
	errStr := "unexpected format of import string (%q), expected ELBID_PROTOCOL_LBPORT_INSTANCEPORT: %s"
	if len(importParts) != 4 {
		return nil, fmt.Errorf(errStr, importStr, "invalid parts")
	}

	elbID := importParts[0]
	protocol := importParts[1]
	lbPort := importParts[2]
	instancePort := importParts[3]

	if elbID == "" {
		return nil, fmt.Errorf(errStr, importStr, "elb id must be required")
	}

	if protocol != "TCP" &&
		protocol != "UDP" &&
		protocol != "HTTP" &&
		protocol != "HTTPS" {
		return nil, fmt.Errorf(errStr, importStr, "protocol must be TCP/UDP/HTTP/HTTPS")
	}

	if _, err := strconv.Atoi(lbPort); err != nil {
		return nil, fmt.Errorf(errStr, importStr, "invalid lb port")
	}
	if _, err := strconv.Atoi(instancePort); err != nil {
		return nil, fmt.Errorf(errStr, importStr, "invalid instance port")
	}
	return importParts, nil
}

func populateELBFromImport(d *schema.ResourceData, importParts []string) error {
	elbID := importParts[0]
	protocol := importParts[1]
	lbPort := importParts[2]
	instancePort := importParts[3]

	if err := d.Set("protocol", protocol); err != nil {
		return err
	}

	p, err := strconv.Atoi(lbPort)
	if err != nil {
		return err
	}

	if err := d.Set("lb_port", p); err != nil {
		return err
	}

	p, err = strconv.Atoi(instancePort)
	if err != nil {
		return err
	}

	if err := d.Set("instance_port", p); err != nil {
		return err
	}

	if err := d.Set("elb_id", elbID); err != nil {
		return err
	}
	return nil
}

func getELBID(d *schema.ResourceData) string {
	return strings.Split(d.Id(), "_")[0]
}
