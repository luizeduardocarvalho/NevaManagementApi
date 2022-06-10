namespace NevaManagement.Infrastructure.Services;

public class ProductUsageService : IProductUsageService
{
    private readonly IProductUsageRepository repository;

    public ProductUsageService(IProductUsageRepository repository)
    {
        this.repository = repository;
    }

    public async Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId)
    {
        try
        {
            return await this.repository.GetLastUsesByResearcher(researcherId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the last uses by researcher.");
        }
    }

    public async Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId)
    {
        try
        {
            return await this.repository.GetLastUsedProductByResearcher(researcherId);
        }
        catch
        {
            throw new Exception("An error occurred while getting last used product by researcher.");
        }
    }
    public async Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId)
    {
        try
        {
            return await this.repository.GetLastUsesByProduct(productId);
        }
        catch
        {
            throw new Exception("An error occurred while getting last uses by product.");
        }
    }
}
