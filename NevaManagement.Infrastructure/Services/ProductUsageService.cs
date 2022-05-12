namespace NevaManagement.Infrastructure.Services;

public class ProductUsageService : IProductUsageService
{
    private readonly IProductUsageRepository repository;
    private readonly IProductRepository productRepository;

    public ProductUsageService(IProductUsageRepository repository, IProductRepository productRepository)
    {
        this.repository = repository;
        this.productRepository = productRepository;
    }

    public async Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId)
    {
        return await this.repository.GetLastUsesByResearcher(researcherId);
    }

    public async Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId)
    {
        try
        {
            return await this.repository.GetLastUsedProductByResearcher(researcherId);
        }
        catch
        {
            throw;
        }
    }
    public async Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId)
    {
        return await this.repository.GetLastUsesByProduct(productId);
    }
}
