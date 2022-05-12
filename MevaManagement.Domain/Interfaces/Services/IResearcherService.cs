namespace NevaManagement.Domain.Interfaces.Services;

public interface IResearcherService
{
    Task<Researcher> GetAll();
    Task<bool> Create(CreateResearcherDto researcherDto);
    Task<IList<GetSimpleResearcherDto>> GetResearchers();
}
