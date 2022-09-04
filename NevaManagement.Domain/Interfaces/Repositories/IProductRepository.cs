namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IProductRepository : IBaseRepository<Product>
{
    Task<IEnumerable<GetProductDto>> GetAll();
    Task<GetProductDto> GetProductById(long id);
    Task<GetDetailedProductDto> GetDetailedProductById(long id);
    Task<IList<GetDetailedProductDto>> GetLowInStockProducts();
}
