namespace NevaManagement.Domain.Interfaces.Services;

public interface IResearcherService
{
    Task<IList<GetSimpleResearcherDto>> GetResearchers(long laboratoryId);
    Task<GetDetailedResearcherDto> GetByEmailAndPassword(string email, string password);
    Task<bool> Create(CreateResearcherDto researcher);
    Task Save();
    Task<Researcher> GetById(long id);
}
