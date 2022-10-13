namespace NevaManagement.Infrastructure.Repositories;

public class EquipmentUsageRepository : BaseRepository<EquipmentUsage>, IEquipmentUsageRepository
{
    private readonly NevaManagementDbContext context;

    public EquipmentUsageRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    async Task<IList<GetEquipmentUsageDto>> IEquipmentUsageRepository.GetEquipmentUsageByEquipment(long equipmentId)
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

    async Task<IList<GetEquipmentUsageDto>> IEquipmentUsageRepository.GetEquipmentUsageByDay(long equipmentId, DateTimeOffset startDate)
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

    async Task<IList<GetEquipmentUsageDto>> IEquipmentUsageRepository.GetEquipmentUsage(long equipmentid, int page)
    {
        return await this.context.EquipmentUsages
            .Where(x => x.EquipmentId == equipmentid)
            .OrderByDescending(x => x.StartDate)
            .Skip((page - 1) * 15)
            .Take(15)
            .Select(x => new GetEquipmentUsageDto
            {
                StartDate = x.StartDate,
                EndDate = x.EndDate,
                Researcher = new GetSimpleResearcherDto
                {
                    Id = x.Researcher.Id,
                    Name = x.Researcher.Name,
                }
            })
            .ToListAsync()
            .ConfigureAwait(false);
    }
}
