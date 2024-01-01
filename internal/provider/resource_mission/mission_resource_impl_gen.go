
package resource_mission

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
var _ resource.Resource = &MissionResource{}
var _ resource.ResourceWithImportState = &MissionResource{}

func NewMissionResource() resource.Resource {
	return &MissionResource{}
}

// MissionResource defines the resource implementation.
type MissionResource struct {
	client *clients.ClientWithResponses
}

func (r *MissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Mission"
}

func (r *MissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = MissionResourceSchema(ctx)
}

func (r *MissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MissionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := clients.CreateMissionJSONRequestBody{}
	js, err := json.Marshal(data)
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Mission, got error: %s", err))
	    return
	}

	if err := json.Unmarshal(js, &body); err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Mission, got error: %s", err))
	    return
	}

	clientResp, err := r.client.CreateMissionWithResponse(ctx, body)
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Mission, got error: %s", err))
	    return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create Mission, got status: %d", clientResp.StatusCode()))
	    return
	}

	// For the purposes of this Mission code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringrValue(clientResp.JSON200.Id.String())

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MissionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	
	clientResp, err := r.client.ReadMissionWithResponse(ctx, data.Id)
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Mission, got error: %s", err))
	    return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to create Mission, got status: %d", clientResp.StatusCode()))
	    return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MissionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}


	body := clients.UpdateMissionJSONRequestBody{}
	js, err := json.Marshal(data)
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Mission, got error: %s", err))
	    return
	}

	if err := json.Unmarshal(js, &body); err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Mission, got error: %s", err))
	    return
	}

	clientResp, err := r.client.UpdateMissionWithResponse(ctx, data.Id.String(), body)
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Mission, got error: %s", err))
	    return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to update Mission, got status: %d", clientResp.StatusCode()))
	    return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MissionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientResp, err := r.client.DeleteMissionWithResponse(ctx, data.Id.String())
	if err != nil {
	    resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Mission, got error: %s", err))
	    return
	}

	if clientResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Server Error", fmt.Sprintf("Unable to delete Mission, got status: %d", clientResp.StatusCode()))
	    return
	}
}

func (r *MissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}