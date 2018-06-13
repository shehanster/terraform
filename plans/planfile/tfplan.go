package planfile

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/plans/internal/planproto"
	"github.com/hashicorp/terraform/version"
)

const tfplanFormatVersion = 3

// ---------------------------------------------------------------------------
// This file deals with the internal structure of the "tfplan" sub-file within
// the plan file format. It's all private API, wrapped by methods defined
// elsewhere. This is the only file that should import the
// ../internal/planproto package, which contains the ugly stubs generated
// by the protobuf compiler.
// ---------------------------------------------------------------------------

// readTFPlan reads a protobuf-encoded description from the plan portion of
// a plan file, which is stored in a special file in the archive called
// "tfplan".
func readTFPlan(r io.Reader) (*plans.Plan, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var rawPlan planproto.Plan
	err = proto.Unmarshal(src, &rawPlan)
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}

	if rawPlan.Version != tfplanFormatVersion {
		return nil, fmt.Errorf("unsupported plan file format version %d; only version %d is supported", rawPlan.Version, tfplanFormatVersion)
	}

	if rawPlan.TerraformVersion != version.String() {
		return nil, fmt.Errorf("plan file was created by Terraform %s, but this is %s; plan files cannot be transferred between different Terraform versions", rawPlan.TerraformVersion, version.String())
	}

	// TODO: Populate the rest of this!
	plan := &plans.Plan{}

	return plan, nil
}

// writeTFPlan serializes the given plan into the protobuf-based format used
// for the "tfplan" portion of a plan file.
func writeTFPlan(plan *plans.Plan, w io.Writer) error {
	rawPlan := &planproto.Plan{
		Version:          tfplanFormatVersion,
		TerraformVersion: version.String(),
	}
	// TODO: Populate the rest of that!

	src, err := proto.Marshal(rawPlan)
	if err != nil {
		return fmt.Errorf("serialization error: %s", err)
	}

	_, err = w.Write(src)
	if err != nil {
		return fmt.Errorf("failed to write plan to plan file: %s", err)
	}

	return nil
}
