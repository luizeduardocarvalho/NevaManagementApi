namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IEquipmentRepository : IBaseRepository<Equipment>
{
    Task<IList<GetSimpleEquipmentDto>> GetEquipments();
    Task<GetDetailedEquipmentDto> GetDetailedEquipment(long id);
}
