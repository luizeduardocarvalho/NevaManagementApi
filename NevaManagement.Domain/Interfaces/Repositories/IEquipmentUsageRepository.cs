namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IEquipmentUsageRepository : IBaseRepository<EquipmentUsage>
{
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByEquipment(long equipmentId);
}
