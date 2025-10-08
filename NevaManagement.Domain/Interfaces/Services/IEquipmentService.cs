namespace NevaManagement.Domain.Interfaces.Services;

public interface IEquipmentService
{
    Task<IList<GetSimpleEquipmentDto>> GetEquipments(long laboratoryId);
    Task<bool> AddEquipment(AddEquipmentDto equipmentDto);
    Task<bool> EditEquipment(EditEquipmentDto equipmentDto);
    Task<GetDetailedEquipmentDto> GetDetailedEquipment(long id, long laboratoryId);
}
