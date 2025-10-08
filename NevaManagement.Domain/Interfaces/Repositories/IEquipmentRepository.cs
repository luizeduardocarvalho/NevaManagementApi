namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IEquipmentRepository : IBaseRepository<Equipment>
{
    Task<IList<GetSimpleEquipmentDto>> GetEquipments(long laboratoryId);
    Task<GetDetailedEquipmentDto> GetDetailedEquipment(long id, long laboratoryId);
}
