

package provider_launchpad

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	
	"github.com/orbit-ops/terraform-provider-launchpad/internal/provider/resource_mission"
	"github.com/orbit-ops/terraform-provider-launchpad/internal/provider/resource_rocket"
	
	"github.com/orbit-ops/terraform-provider-launchpad/internal/provider/datasource_mission"
	"github.com/orbit-ops/terraform-provider-launchpad/internal/provider/datasource_rocket"

	"github.com/orbit-ops/terraform-provider-launchpad/internal/clients"
)

// Ensure LaunchpadProvider satisfies various provider interfaces.
var _ provider.Provider = &LaunchpadProvider{}

// LaunchpadProvider defines the provider implementation.
type LaunchpadProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *LaunchpadProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "Launchpad"
	resp.Version = p.version
}

func (p *LaunchpadProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = LaunchpadProviderSchema(ctx)
}

func (p *LaunchpadProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LaunchpadModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := clients.NewClientWithResponses(data.Endpoint.String())
	if err != nil {
		resp.Diagnostics.AddError("Provider error", fmt.Sprintf("got error: %v", err))
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LaunchpadProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
            resource_mission.NewMissionResource,
            resource_rocket.NewRocketResource,
	}
}

func (p *LaunchpadProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
            datasource_mission.NewMissionDataSource,
            datasource_rocket.NewRocketDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LaunchpadProvider{
			version: version,
		}
	}
}