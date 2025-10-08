using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Domain.Interfaces.Services;

public interface ILaboratoryService
{
    Task<IList<GetLaboratoryDto>> GetAllLaboratories();
    Task<IList<GetSimpleLaboratoryDto>> GetSimpleLaboratories();
    Task<GetLaboratoryDto> GetLaboratoryById(long id);
    Task<bool> CreateLaboratory(CreateLaboratoryDto createLaboratoryDto);
    Task<bool> UpdateLaboratory(EditLaboratoryDto editLaboratoryDto);
    Task<bool> DeleteLaboratory(long id);
    Task<bool> LaboratoryExists(long id);
}
