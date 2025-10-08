namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IResearcherRepository : IBaseRepository<Researcher>
{
    Task<IEnumerable<Researcher>> GetAll(long laboratoryId);
    Task<IList<GetSimpleResearcherDto>> GetResearchers(long laboratoryId);
    Task<GetDetailedResearcherDto> GetByEmailAndPassword(string email, string password);
}
