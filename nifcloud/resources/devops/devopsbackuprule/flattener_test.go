package devopsbackuprule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops"
	"github.com/nifcloud/nifcloud-sdk-go/service/devops/types"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})
	rd.SetId("test_name")

	wantRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{
		"name":        "test_name",
		"instance_id": "test_id",
		"description": "test_description",
		"backup_time": "test_time",
	})
	wantRd.SetId("test_name")

	wantNotFoundRd := schema.TestResourceDataRaw(t, newSchema(), map[string]interface{}{})

	type args struct {
		res *devops.GetBackupRuleOutput
		d   *schema.ResourceData
	}
	tests := []struct {
		name string
		args args
		want *schema.ResourceData
	}{
		{
			name: "flattens the response",
			args: args{
				d: rd,
				res: &devops.GetBackupRuleOutput{
					BackupRule: &types.BackupRule{
						BackupRuleName: nifcloud.String("test_name"),
						InstanceId:     nifcloud.String("test_id"),
						Description:    nifcloud.String("test_description"),
						BackupTime:     nifcloud.String("test_time"),
					},
				},
			},
			want: wantRd,
		},
		{
			name: "flattens the response even when the resource has been removed externally",
			args: args{
				d:   wantNotFoundRd,
				res: nil,
			},
			want: wantNotFoundRd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flatten(tt.args.d, tt.args.res)
			assert.NoError(t, err)

			wantState := tt.want.State()
			if wantState == nil {
				tt.want.SetId("some")
				wantState = tt.want.State()
			}

			gotState := tt.args.d.State()
			if gotState == nil {
				tt.args.d.SetId("some")
				gotState = tt.args.d.State()
			}

			assert.Equal(t, wantState.Attributes, gotState.Attributes)
		})
	}
}
