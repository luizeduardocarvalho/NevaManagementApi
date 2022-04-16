using NevaManagement.Domain.Dtos.Researcher;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IResearcherRepository
    {
        Task<IEnumerable<Researcher>> GetAll();
        Task<bool> Create(Researcher researcher);
        Task<Researcher> GetEntityById(long id);
        Task<bool> SaveChanges();
        Task<IList<GetSimpleResearcherDto>> GetResearchers();
    }
}
