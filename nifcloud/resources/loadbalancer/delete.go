package loadbalancer

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/terraform-provider-nifcloud/nifcloud/client"
)

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := meta.(*client.Client).Computing
	lbns := []computing.RequestLoadBalancerNames{
		{
			LoadBalancerName: nifcloud.String(d.Get("load_balancer_name").(string)),
		},
	}
	req := svc.DescribeLoadBalancersRequest(&computing.DescribeLoadBalancersInput{
		LoadBalancerNames: lbns,
	})

	res, err := req.Send(ctx)
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "Client.InvalidParameterNotFound.LoadBalancerName" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed reading: %s", err))
	}

	for _, lbd := range res.DescribeLoadBalancersOutput.DescribeLoadBalancersResult.LoadBalancerDescriptions {
		n := lbd.ListenerDescriptions[0].Listener
		req := svc.DeleteLoadBalancerRequest(&computing.DeleteLoadBalancerInput{
			InstancePort:     nifcloud.Int64(int64(int(*n.InstancePort))),
			LoadBalancerName: nifcloud.String(d.Id()),
			LoadBalancerPort: nifcloud.Int64(int64(int(*n.LoadBalancerPort))),
		})
		_, err := req.Send(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed deleting load_balancer: %s", err))
		}
	}
	d.SetId("")
	return nil
}
