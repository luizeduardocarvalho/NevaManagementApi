namespace NevaManagement.Domain.Interfaces.Services;

public interface IResearcherService
{
    Task<Researcher> GetAll();
    Task<bool> Create(CreateResearcherDto researcherDto);
    Task<IList<GetSimpleResearcherDto>> GetResearchers();
    Task<Researcher> GetByEmailAndPassword(string email, string password);
}
