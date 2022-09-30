namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IProductRepository : IBaseRepository<Product>
{
    Task<IEnumerable<GetProductDto>> GetAll(int page);
    Task<GetProductDto> GetProductById(long id);
    Task<GetDetailedProductDto> GetDetailedProductById(long id);
}
