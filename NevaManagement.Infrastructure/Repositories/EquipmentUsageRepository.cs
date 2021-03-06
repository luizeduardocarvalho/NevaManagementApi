using NevaManagement.Domain.Dtos.EquipmentUsage;

namespace NevaManagement.Infrastructure.Repositories;

public class EquipmentUsageRepository : BaseRepository<EquipmentUsage>, IEquipmentUsageRepository
{
    private readonly NevaManagementDbContext context;

    public EquipmentUsageRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByEquipment(long equipmentId)
    {
        return await this.context.EquipmentUsages
            .Where(equipmentUsage => equipmentUsage.EquipmentId == equipmentId)
            .Where(equipmentUsage => equipmentUsage.EndDate >= DateTimeOffset.UtcNow)
            .Select(equipmentUsage => new GetEquipmentUsageDto()
            {
                Id = equipmentUsage.Id,
                Researcher = new GetSimpleResearcherDto()
                {
                    Id = equipmentUsage.Researcher.Id,
                    Name = equipmentUsage.Researcher.Name
                },
                StartDate = equipmentUsage.StartDate,
                EndDate = equipmentUsage.EndDate,
            })
            .ToListAsync();
    }

    public async Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByDay(long equipmentId, DateTimeOffset startDate)
    {
        return await this.context.EquipmentUsages
            .Where(equipmentUsage => equipmentUsage.EquipmentId == equipmentId)
            .Where(equipmentUsage => equipmentUsage.StartDate.Day == startDate.Day)
            .Select(equipmentUsage => new GetEquipmentUsageDto
            {
                StartDate = equipmentUsage.StartDate,
                EndDate = equipmentUsage.EndDate
            })
            .ToListAsync();
    }
}
