package planfile

import (
	"bytes"
	"testing"

	"github.com/go-test/deep"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/plans"
)

func TestTFPlanRoundTrip(t *testing.T) {
	plan := &plans.Plan{
		VariableValues: map[string]cty.Value{
			"foo": cty.StringVal("foo value"),
		},
		Changes: &plans.Changes{
			RootOutputs: map[string]*plans.OutputChange{
				"bar": {
					Change: plans.Change{
						Action: plans.Create,
						Before: cty.NullVal(cty.String),
						After:  cty.StringVal("bar value"),
					},
					Sensitive: false,
				},
				"baz": {
					Change: plans.Change{
						Action: plans.NoOp,
						Before: cty.StringVal("baz value"),
						After:  cty.StringVal("baz value"),
					},
					Sensitive: false,
				},
				"secret": {
					Change: plans.Change{
						Action: plans.Update,
						Before: cty.StringVal("old secret value"),
						After:  cty.StringVal("new secret value"),
					},
					Sensitive: true,
				},
			},
			Resources: []*plans.ResourceInstanceChange{
				{
					Addr: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_thing",
						Name: "woot",
					}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
					Change: plans.Change{
						Action: plans.Replace,
						Before: cty.ObjectVal(map[string]cty.Value{
							"id": cty.StringVal("foo-bar-baz"),
						}),
						After: cty.ObjectVal(map[string]cty.Value{
							"id": cty.UnknownVal(cty.String),
						}),
					},
				},
				{
					Addr: addrs.Resource{
						Mode: addrs.ManagedResourceMode,
						Type: "test_thing",
						Name: "woot",
					}.Instance(addrs.IntKey(0)).Absolute(addrs.RootModuleInstance),
					DeposedKey: "foodface",
					Change: plans.Change{
						Action: plans.Delete,
						Before: cty.ObjectVal(map[string]cty.Value{
							"id": cty.StringVal("bar-baz-foo"),
						}),
						After: cty.NullVal(cty.String),
					},
				},
			},
		},
		ProviderSHA256s: map[string][]byte{
			"test": []byte{
				0xba, 0x5e, 0x1e, 0x55, 0xb0, 0x1d, 0xfa, 0xce,
				0xef, 0xfe, 0xc7, 0xed, 0x1a, 0xbe, 0x11, 0xed,
				0x5c, 0xa1, 0xab, 0x1e, 0xda, 0x7a, 0xba, 0x5e,
				0x70, 0x7a, 0x11, 0xed, 0xb0, 0x07, 0xab, 0x1e,
			},
		},
	}

	var buf bytes.Buffer
	err := writeTFPlan(plan, &buf)
	if err != nil {
		t.Fatal(err)
	}

	newPlan, err := readTFPlan(&buf)
	if err != nil {
		t.Fatal(err)
	}

	{
		oldDepth := deep.MaxDepth
		oldCompare := deep.CompareUnexportedFields
		deep.MaxDepth = 20
		deep.CompareUnexportedFields = true
		defer func() {
			deep.MaxDepth = oldDepth
			deep.CompareUnexportedFields = oldCompare
		}()
	}
	for _, problem := range deep.Equal(newPlan, plan) {
		t.Error(problem)
	}
}
