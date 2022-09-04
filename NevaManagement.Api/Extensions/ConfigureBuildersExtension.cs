using NevaManagement.Domain.Interfaces.Builders;
using NevaManagement.Infrastructure.Builders;

namespace NevaManagement.Api.Extensions;

public static class ConfigureBuildersExtension
{
    public static void ConfigureBuilders(this IServiceCollection services)
    {
        services.AddTransient<IBuilder<AddOrganismDto, Organism>, OrganismBuilder>();
    }
}
