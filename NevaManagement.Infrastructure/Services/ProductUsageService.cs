using NevaManagement.Domain.Dtos.ProductUsage;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
    public class ProductUsageService : IProductUsageService
    {
        private readonly IProductUsageRepository repository;

        public ProductUsageService(IProductUsageRepository repository)
        {
            this.repository = repository;
        }

        public async Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId)
        {
            return await this.repository.GetLastUsesByResearcher(researcherId);
        }

        public async Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId)
        {
            return await this.repository.GetLastUsedProductByResearcher(researcherId);
        }
    }
}
