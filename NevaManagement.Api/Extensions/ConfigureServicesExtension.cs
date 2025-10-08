namespace NevaManagement.Api.Extensions;

public static class ConfigureServicesExtension
{
    public static void ConfigureServices(this IServiceCollection services)
    {
        services.AddTransient<IResearcherService, ResearcherService>();
        services.AddTransient<IProductService, ProductService>();
        services.AddTransient<IProductUsageRepository, ProductUsageRepository>();
        services.AddTransient<ILocationService, LocationService>();
        services.AddTransient<IOrganismService, OrganismService>();
        services.AddTransient<IContainerService, ContainerService>();
        services.AddTransient<ITokenService, TokenService>();
        services.AddTransient<IEncryptService, EncryptService>();
        services.AddTransient<IAuthService, AuthService>();
        services.AddTransient<IArticleService, ArticleService>();
        services.AddTransient<IEquipmentService, EquipmentService>();
        services.AddTransient<IEquipmentUsageService, EquipmentUsageService>();
        services.AddTransient<ILaboratoryService, LaboratoryService>();
        services.AddTransient<ILaboratoryInvitationService, LaboratoryInvitationService>();
    }
}
