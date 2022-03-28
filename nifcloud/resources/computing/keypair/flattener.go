package keypair

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeKeyPairsOutput) error {
	if res == nil || len(res.KeySet) == 0 {
		d.SetId("")
		return nil
	}

	keyPair := res.KeySet[0]

	if nifcloud.ToString(keyPair.KeyName) != d.Id() {
		return fmt.Errorf("unable to find key pair within: %#v", res.KeySet)
	}

	if err := d.Set("key_name", keyPair.KeyName); err != nil {
		return err
	}

	if err := d.Set("fingerprint", keyPair.KeyFingerprint); err != nil {
		return err
	}

	if err := d.Set("description", keyPair.Description); err != nil {
		return err
	}
	return nil
}
