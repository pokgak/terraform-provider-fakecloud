package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fakecloud "github.com/pokgak/fakecloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &virtualMachinesDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualMachinesDataSource{}
)

func NewVirtualMachinesDataSource() datasource.DataSource {
	return &virtualMachinesDataSource{}
}

type virtualMachinesDataSource struct {
	client *fakecloud.Client
}

// virtualMachinesDataSourceModel maps the data source schema data.
type virtualMachinesDataSourceModel struct {
	VirtualMachines []virtualMachineModel `tfsdk:"virtual_machines"`
}

// virtualMachineModel maps coffees schema data.
type virtualMachineModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	InstanceType types.String `tfsdk:"instance_type"`
}

func (d *virtualMachinesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machines"
}

// Configure adds the provider configured client to the data source.
func (d *virtualMachinesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fakecloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fakecloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Schema defines the schema for the data source.
func (d *virtualMachinesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "virtual_machines": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.Int64Attribute{
                            Computed: true,
                        },
                        "name": schema.StringAttribute{
                            Computed: true,
                        },
                        "instance_type": schema.StringAttribute{
                            Computed: true,
                        },
                    },
                },
            },
        },
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *virtualMachinesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualMachinesDataSourceModel

	vms, err := d.client.GetVMs()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Fakecloud VMs",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, vm := range vms {
		vmState := virtualMachineModel{
			ID:           types.Int64Value(int64(vm.ID)),
			Name:         types.StringValue(vm.Name),
			InstanceType: types.StringValue(vm.InstanceType),
		}

		state.VirtualMachines = append(state.VirtualMachines, vmState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
