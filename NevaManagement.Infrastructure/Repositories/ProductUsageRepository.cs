namespace NevaManagement.Infrastructure.Repositories;

public class ProductUsageRepository : BaseRepository<ProductUsage>, IProductUsageRepository
{
    private readonly NevaManagementDbContext context;

    public ProductUsageRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId)
    {
        return await this.context.ProductUsages
            .Where(x => x.Researcher.Id == researcherId)
            .Select(x => new GetLastUsesByResearcherDto
            {
                Quantity = x.Quantity,
                Name = x.Product.Name,
                Unit = x.Product.Unit,
                UsageDate = x.UsageDate
            })
            .ToListAsync();
    }

    public async Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId)
    {
        return await this.context.ProductUsages
            .Where(x => x.Product.Id == productId)
            .Select(x => new GetLastUseByProductDto
            {
                Quantity = x.Quantity,
                ResearcherName = x.Researcher.Name,
                Unit = x.Product.Unit,
                UseDate = x.UsageDate
            })
            .ToListAsync();
    }

    public async Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId)
    {
        return await this.context.ProductUsages
            .Where(x => x.Id == researcherId)
            .OrderByDescending(x => x.UsageDate)
            .Select(x => new GetLastUsedProductDto
            {
                ProductId = x.Product.Id,
                LocationName = x.Product.Location.Name,
                ProductName = x.Product.Name,
                Quantity = x.Quantity,
                Unit = x.Product.Unit
            })
            .FirstOrDefaultAsync();
    }

    public async Task<IList<GetLastUseByProductDto>> GetLastThreeMonthsUsesForAllProducts()
    {
        var todayDate = DateTime.Today;
        return await this.context.ProductUsages
            .Where(x => x.UsageDate.Year == todayDate.Year && 
                todayDate.Month - x.UsageDate.Month <= 3)
            .Select(x => new GetLastUseByProductDto
            {
                Product = new GetDetailedProductDto
                {
                    Id = x.Product.Id,
                    Name = x.Product.Name,
                    Quantity = x.Product.Quantity,
                    Unit = x.Product.Unit,
                },
                Quantity = x.Quantity
            })
            .ToListAsync();
    }
}
