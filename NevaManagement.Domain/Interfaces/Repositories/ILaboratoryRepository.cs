using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Domain.Interfaces.Repositories;

public interface ILaboratoryRepository : IBaseRepository<Laboratory>
{
    Task<IList<GetLaboratoryDto>> GetAllLaboratories();
    Task<IList<GetSimpleLaboratoryDto>> GetSimpleLaboratories();
    Task<GetLaboratoryDto> GetLaboratoryById(long id);
    Task<bool> LaboratoryExists(long id);
}
