namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IEquipmentUsageRepository : IBaseRepository<EquipmentUsage>
{
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByEquipment(long equipmentId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByDay(long equipmentId, DateTimeOffset startDate);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsage(long equipmentid, int page);
}
