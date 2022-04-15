using NevaManagement.Domain.Dtos.Organism;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IOrganismRepository
    {
        Task<GetDetailedOrganismDto> GetById(long id);
        Task<Organism> GetEntityById(long id);
        Task<IList<GetOrganismDto>> GetOrganisms();
        Task<bool> Create(Organism organism);
        Task<bool> SaveChanges();
    }
}
