package resource_rocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orbit-ops/terraform-provider-launchpad/internal/clients"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RocketResource{}
var _ resource.ResourceWithImportState = &RocketResource{}

func NewRocketResource() resource.Resource {
	return &RocketResource{}
}

// RocketResource defines the resource implementation.
type RocketResource struct {
	client *clients.ClientWithResponses
}

func (r *RocketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Rocket"
}

func (r *RocketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = RocketResourceSchema(ctx)
}

func (r *RocketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.ClientWithResponses)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clients.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RocketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RocketModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := clients.CreateRocketJSONRequestBody{}
	js, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Rocket, got error: %s", err))
		return
	}

	if err := json.Unmarshal(js, &body); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Rocket, got error: %s", err))
		return
	}

	clientResp, err := r.client.CreateRocketWithResponse(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Rocket, got error: %s", err))
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create Rocket, got status: %d", clientResp.StatusCode()))
		return
	}

	// For the purposes of this Rocket code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(clientResp.JSON200.Id.String())

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RocketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RocketModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.ReadRocketWithResponse(ctx, data.Id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Rocket, got error: %s", err))
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create Rocket, got status: %d", clientResp.StatusCode()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RocketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RocketModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := clients.UpdateRocketJSONRequestBody{}
	js, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Rocket, got error: %s", err))
		return
	}

	if err := json.Unmarshal(js, &body); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Rocket, got error: %s", err))
		return
	}

	clientResp, err := r.client.UpdateRocketWithResponse(ctx, data.Id.String(), body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Rocket, got error: %s", err))
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to update Rocket, got status: %d", clientResp.StatusCode()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RocketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RocketModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.DeleteRocketWithResponse(ctx, data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Rocket, got error: %s", err))
		return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to delete Rocket, got status: %d", clientResp.StatusCode()))
		return
	}
}

func (r *RocketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
