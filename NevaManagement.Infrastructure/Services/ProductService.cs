namespace NevaManagement.Infrastructure.Services;

public class ProductService : IProductService
{
    private readonly IProductRepository repository;
    private readonly ILocationRepository locationRepository;
    private readonly IResearcherRepository researcherRepository;
    private readonly IProductUsageRepository productUsageRepository;

    public ProductService(
        IProductRepository repository,
        ILocationRepository locationRepository,
        IResearcherRepository researcherRepository,
        IProductUsageRepository productUsageRepository)
    {
        this.repository = repository;
        this.locationRepository = locationRepository;
        this.researcherRepository = researcherRepository;
        this.productUsageRepository = productUsageRepository;
    }

    public async Task<IEnumerable<GetProductDto>> GetAll(int page)
    {
        try
        {
            return await this.repository.GetAll(page);
        }
        catch (Exception ex)
        {
            throw new Exception("An error occurred while getting all products.");
        }
    }

    public async Task<GetDetailedProductDto> GetDetailedProductById(long id)
    {
        try
        {
            return await this.repository.GetDetailedProductById(id);
        }
        catch (Exception ex)
        {
            throw new Exception("An error occurred while getting a detailed product.");
        }
    }

    public async Task<GetProductDto> GetById(long id)
    {
        return await this.repository.GetProductById(id);
    }

    public async Task<bool> Create(CreateProductDto productDto)
    {
        try
        {
            var location = await this.locationRepository.GetById(productDto.LocationId.Value);

            if (location is null)
            {
                return false;
            }

            var product = new Product
            {
                Description = productDto.Description,
                Location = location,
                Name = productDto.Name,
                Quantity = productDto.Quantity,
                Unit = productDto.Unit,
            };

            if (productDto.ParsedExpirationDate != null)
            {
                product.ExpirationDate = productDto.ParsedExpirationDate.Value.UtcDateTime;
            }

            await this.repository.Insert(product);
            return await this.repository.SaveChanges();

        }
        catch (DbUpdateException ex)
        {
            throw new DbUpdateException("An error occurred while creating the new product.");
        }
    }

    public async Task<bool> AddQuantityToProduct(AddQuantityToProductDto addQuantityToProductDto)
    {
        var product = await this.repository.GetById(addQuantityToProductDto.ProductId);

        if (product is null)
        {
            return false;
        }

        try
        {
            product.Quantity += addQuantityToProductDto.Quantity;
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while adding quantity to product.");
        }
    }

    public async Task<bool> UseProduct(UseProductDto useProductDto)
    {
        var product = await this.repository.GetById(useProductDto.ProductId);

        if (product is null)
        {
            return false;
        }

        try
        {
            product.Quantity -= useProductDto.Quantity;
            var result = await this.repository.SaveChanges();

            if (!result)
            {
                return false;
            }

            var researcher = await this.researcherRepository.GetById(useProductDto.ResearcherId);

            var productUsage = new ProductUsage
            {
                Researcher = researcher,
                Product = product,
                UsageDate = DateTimeOffset.Now.UtcDateTime,
                Description = useProductDto.Description,
                Quantity = useProductDto.Quantity
            };

            await this.productUsageRepository.Insert(productUsage);
            return await this.productUsageRepository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while using the product.");
        }
    }

    public async Task<bool> EditProduct(EditProductDto editProductDto)
    {
        var product = await this.repository.GetById(editProductDto.Id.Value);

        if (product is null)
        {
            return false;
        }

        product.Name = editProductDto.Name;
        product.Description = editProductDto.Description;
        product.Formula = editProductDto.Formula;

        try
        {
            if (editProductDto.LocationId is not null)
            {
                var location = await this.locationRepository.GetById(editProductDto.LocationId.Value);
                product.Location = location;
            }

            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while editing the product.");
        }
    }

    public async Task<IList<GetDetailedProductDto>> GetLowInStockProducts()
    {
        var lastUses = await this.productUsageRepository.GetLastThreeMonthsUsesForAllProducts();

        var totalUsedPerProduct = lastUses
            .GroupBy(x => x.Product.Id)
            .Select(x => new
            {
                TotalQuantityUsed = x.Sum(p => p.Quantity),
                ProductId = x.Key,
                ProductName = x.First().Product.Name,
                QuantityInStock = x.First().Product.Quantity,
                Unit = x.First().Product.Unit,
            })
            .Where(x => x.TotalQuantityUsed >= x.QuantityInStock);

        var lowInStockProducts = new List<GetDetailedProductDto>();

        foreach (var product in totalUsedPerProduct)
        {
            var lowInStockProduct = new GetDetailedProductDto
            {
                Id = product.ProductId,
                Name = product.ProductName,
                Quantity = product.QuantityInStock,
                Unit = product.Unit,
                QuantityUsedInTheLastThreeMonths = product.TotalQuantityUsed
            };

            lowInStockProducts.Add(lowInStockProduct);
        }

        return lowInStockProducts;
    }
}
