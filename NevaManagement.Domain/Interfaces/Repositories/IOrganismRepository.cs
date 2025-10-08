namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IOrganismRepository : IBaseRepository<Organism>
{
    Task<GetDetailedOrganismDto> GetDetailedOrganismById(long id, long laboratoryId);
    Task<IList<GetOrganismDto>> GetOrganisms(long laboratoryId);
}
