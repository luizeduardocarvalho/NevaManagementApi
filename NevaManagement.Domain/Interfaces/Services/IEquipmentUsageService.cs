namespace NevaManagement.Domain.Interfaces.Services;

public interface IEquipmentUsageService
{
    Task<IList<EquipmentUsageCalendarDto>> GetEquipmentUsageCalendar(long equipmentId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageHistory(long id);
    Task<bool> UseEquipment(UseEquipmentDto equipmentDto);
}
