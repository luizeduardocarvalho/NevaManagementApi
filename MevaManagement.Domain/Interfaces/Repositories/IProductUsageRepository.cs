using NevaManagement.Domain.Dtos.ProductUsage;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IProductUsageRepository
    {
        Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId);
    }
}
