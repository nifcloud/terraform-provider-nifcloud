package devopsrunnerregistration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func createRunnerRegistration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).DevOpsRunner

	input := expandRegisterRunnerInput(d)

	_, err := svc.RegisterRunner(ctx, input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to register the DevOps Runner to a GitLab instance: %s", err))
	}

	err = waitUntilRunnerRunning(ctx, d, svc)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to wait for the DevOps Runner to become ready: %s", err))
	}

	regRes, err := svc.ListRunnerRegistrations(ctx, expandListRunnerRegistrationsInput(d))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read a list of DevOps Runner registrations: %s", err))
	}

	d.SetId(nifcloud.ToString(regRes.Registrations[len(regRes.Registrations)-1].RegistrationId))

	return readRunnerRegistration(ctx, d, meta)
}
