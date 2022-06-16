namespace NevaManagement.Domain.Interfaces.Services;

public interface IEquipmentUsageService
{
    Task<dynamic> GetEquipmentUsageCalendar(long equipmentId);
    Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageHistory(long id);
    Task<bool> UseEquipment(UseEquipmentDto equipmentDto);
}
