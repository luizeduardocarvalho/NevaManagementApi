namespace NevaManagement.Domain.Interfaces.Services;

public interface IResearcherService
{
    Task<IList<GetSimpleResearcherDto>> GetResearchers();
    Task<Researcher> GetByEmailAndPassword(string email, string password);
    Task<bool> Create(CreateResearcherDto researcher);
}
