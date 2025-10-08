namespace NevaManagement.Api.Extensions;

public static class ConfigureRepositoriesExtension
{
    public static void ConfigureRepositories(this IServiceCollection services)
    {
        services.AddTransient<IResearcherRepository, ResearcherRepository>();
        services.AddTransient<IProductRepository, ProductRepository>();
        services.AddTransient<ILocationRepository, LocationRepository>();
        services.AddTransient<IProductUsageService, ProductUsageService>();
        services.AddTransient<ILocationRepository, LocationRepository>();
        services.AddTransient<IOrganismRepository, OrganismRepository>();
        services.AddTransient<IContainerRepository, ContainerRepository>();
        services.AddTransient<IArticleRepository, ArticleRepository>();
        services.AddTransient<IArticleContainerRepository, ArticleContainerRepository>();
        services.AddTransient<IEquipmentRepository, EquipmentRepository>();
        services.AddTransient<IEquipmentUsageRepository, EquipmentUsageRepository>();
        services.AddTransient<ILaboratoryRepository, LaboratoryRepository>();
        services.AddTransient<ILaboratoryInvitationRepository, LaboratoryInvitationRepository>();
    }
}
