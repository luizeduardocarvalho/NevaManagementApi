using NevaManagement.Domain.Dtos.ProductUsage;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IProductUsageRepository
    {
        Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId);
        Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId);
        Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId);
        Task<bool> SaveChanges();
        Task<bool> Create(ProductUsage productUsage);
    }
}
