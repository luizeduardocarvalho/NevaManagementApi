namespace NevaManagement.Domain.Interfaces.Services;

public interface IEquipmentUsageService
{
    Task<IList<EquipmentUsageCalendarDto>> GetEquipmentUsageCalendar(long equipmentId, long laboratoryId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageHistory(long id, long laboratoryId);
    Task<bool> UseEquipment(UseEquipmentDto equipmentDto);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsage(long equipmentid, int page, long laboratoryId);
}
