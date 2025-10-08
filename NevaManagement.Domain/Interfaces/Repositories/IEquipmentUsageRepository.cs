namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IEquipmentUsageRepository : IBaseRepository<EquipmentUsage>
{
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByEquipment(long equipmentId, long laboratoryId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageByDay(long equipmentId, DateTimeOffset startDate, long laboratoryId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsage(long equipmentid, int page, long laboratoryId);
}
