namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IResearcherRepository : IBaseRepository<Researcher>
{
    Task<IEnumerable<Researcher>> GetAll();
    Task<IList<GetSimpleResearcherDto>> GetResearchers();
    Task<GetDetailedResearcherDto> GetByEmailAndPassword(string email, string password);
}
