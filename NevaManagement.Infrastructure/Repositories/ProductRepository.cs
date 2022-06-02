namespace NevaManagement.Infrastructure.Repositories;

public class ProductRepository : BaseRepository<Product>, IProductRepository
{
    private readonly NevaManagementDbContext context;

    public ProductRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IEnumerable<GetProductDto>> GetAll()
    {
        return await this.context.Products
                                    .Select(p => new GetProductDto
                                    {
                                        Id = p.Id,
                                        Name = p.Name
                                    }).ToListAsync();
    }

    public async Task<GetProductDto> GetById(long id)
    {
        return await this.context.Products
            .Where(x => x.Id == id)
            .Select(x => new GetProductDto
            {
                Id = x.Id,
                Name = x.Name
            })
            .SingleOrDefaultAsync();
    }

    public async Task<GetDetailedProductDto> GetDetailedProductById(long id)
    {
        return await this.context.Products
                                    .Where(x => x.Id == id)
                                    .Include(x => x.Location)
                                    .Select(x => new GetDetailedProductDto
                                    {
                                        Id = x.Id,
                                        Location = new GetLocationDto
                                        {
                                            Id = x.Location.Id,
                                            Name = x.Location.Name
                                        },
                                        Name = x.Name,
                                        Quantity = x.Quantity,
                                        Unit = x.Unit,
                                        Description = x.Description,
                                        ExpirationDate = x.ExpirationDate
                                    })
                                    .FirstOrDefaultAsync();
    }

    public async Task<Product> GetEntityById(long id)
    {
        return await this.context.Products.Where(x => x.Id == id).FirstOrDefaultAsync();
    }

    public new async Task<bool> Insert(Product product)
    {
        await this.context.AddAsync(product);

        return await SaveChanges();
    }

    public new async Task<bool> SaveChanges()
    {
        var result = await this.context.SaveChangesAsync();

        if (result > 0)
        {
            return true;
        }

        return false;
    }
}
