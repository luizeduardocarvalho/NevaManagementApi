using NevaManagement.Domain.Dtos.EquipmentUsage;

namespace NevaManagement.Api.Extensions;

public static class ConfigureAutoMapperExtension
{
    public static void ConfigureAutoMapper(this IServiceCollection services)
    {
        var config = new MapperConfiguration(cfg =>
        {
            cfg.CreateMap<AddEquipmentDto, Equipment>();
            cfg.CreateMap<UseEquipmentDto, EquipmentUsage>();
        });

        var mapper = config.CreateMapper();
        services.AddSingleton(mapper);
    }
}
