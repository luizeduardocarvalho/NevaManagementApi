namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IOrganismRepository : IBaseRepository<Organism>
{
    Task<GetDetailedOrganismDto> GetDetailedOrganismById(long id);
    Task<IList<GetOrganismDto>> GetOrganisms();
}
