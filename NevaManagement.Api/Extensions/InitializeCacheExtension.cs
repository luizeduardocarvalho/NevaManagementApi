namespace NevaManagement.Api.Extensions;

public static class InitializeCacheExtension
{
    public static async Task<IApplicationBuilder> InitializeCache(this IApplicationBuilder builder)
    {
        using var scope = builder.ApplicationServices.CreateScope();
        var locationService = scope.ServiceProvider.GetService<ILocationService>();
        _ = await locationService.GetCachedLocations();

        return builder;
    }
}
