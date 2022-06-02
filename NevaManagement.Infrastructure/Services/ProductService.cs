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

    public async Task<IEnumerable<GetProductDto>> GetAll()
    {
        return await this.repository.GetAll();
    }

    public async Task<GetDetailedProductDto> GetDetailedProductById(long id)
    {
        return await this.repository.GetDetailedProductById(id);
    }

    public async Task<GetProductDto> GetById(long id)
    {
        return await this.repository.GetById(id);
    }

    public async Task<bool> Create(CreateProductDto productDto)
    {
        try
        {
            var location = await this.locationRepository.GetEntityById(productDto.LocationId.Value);

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
                ExpirationDate = productDto.ExpirationDate.UtcDateTime
            };

            return await this.repository.Insert(product);

        }
        catch(DbUpdateException ex)
        {
            throw new DbUpdateException("An error occurred while creating the new product.");
        }
    }

    public async Task<bool> AddQuantityToProduct(AddQuantityToProductDto addQuantityToProductDto)
    {
        var product = await this.repository.GetEntityById(addQuantityToProductDto.ProductId);
        var result = false;

        if (product is not null)
        {
            try
            {
                product.Quantity += addQuantityToProductDto.Quantity;
                result = await this.repository.SaveChanges();
            }
            catch
            {
                return false;
            }
        }

        return result;
    }

    public async Task<bool> UseProduct(UseProductDto useProductDto)
    {
        var product = await this.repository.GetEntityById(useProductDto.ProductId);
        var result = false;

        if (product is not null)
        {
            try
            {
                product.Quantity -= useProductDto.Quantity;
                result = await this.repository.SaveChanges();

                var researcher = await this.researcherRepository.GetEntityById(useProductDto.ResearcherId);

                var productUsage = new ProductUsage
                {
                    Researcher = researcher,
                    Product = product,
                    UsageDate = DateTimeOffset.Now,
                    Description = useProductDto.Description,
                    Quantity = useProductDto.Quantity
                };

                result = await this.productUsageRepository.Create(productUsage);
            }
            catch
            {
                return false;
            }
        }

        return result;
    }

    public async Task<bool> EditProduct(EditProductDto editProductDto)
    {
        var result = false;
        var product = await this.repository.GetEntityById(editProductDto.Id.Value);

        if (product is not null)
        {
            product.Name = editProductDto.Name;
            product.Description = editProductDto.Description;

            try
            {
                if (editProductDto.LocationId is not null)
                {
                    var location = await this.locationRepository.GetEntityById(editProductDto.LocationId.Value);
                    product.Location = location;
                }

                result = await this.repository.SaveChanges();
            }
            catch
            {
                throw;
            }

        }

        return result;
    }
}
