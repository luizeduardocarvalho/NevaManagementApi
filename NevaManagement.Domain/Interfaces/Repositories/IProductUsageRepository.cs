namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IProductUsageRepository : IBaseRepository<ProductUsage>
{
    Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId);
    Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId);
    Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId);
}
