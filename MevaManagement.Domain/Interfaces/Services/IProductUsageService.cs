namespace NevaManagement.Domain.Interfaces.Services;

public interface IProductUsageService
{
    Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId);
    Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId);
    Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId);
}
