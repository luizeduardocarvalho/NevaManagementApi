namespace NevaManagement.Infrastructure.Repositories;

public class EquipmentRepository : BaseRepository<Equipment>, IEquipmentRepository
{
    private readonly NevaManagementDbContext context;

    public EquipmentRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IList<GetSimpleEquipmentDto>> GetEquipments()
    {
        return await this.context.Equipments
            .Select(equipment => new GetSimpleEquipmentDto()
            {
                Id = equipment.Id,
                Name = equipment.Name
            })
            .ToListAsync();
    }

    public async Task<GetDetailedEquipmentDto> GetDetailedEquipment(long id)
    {
        return await this.context.Equipments
            .Where(equipment => equipment.Id == id)
            .Select(equipment => new GetDetailedEquipmentDto
            {
                Name = equipment.Name,
                Description = equipment.Description
            })
            .FirstOrDefaultAsync();
    }
}
