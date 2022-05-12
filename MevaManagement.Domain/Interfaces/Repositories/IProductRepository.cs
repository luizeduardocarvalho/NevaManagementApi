namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IProductRepository
{
    Task<IEnumerable<GetProductDto>> GetAll();
    Task<GetProductDto> GetById(long id);
    Task<GetDetailedProductDto> GetDetailedProductById(long id);
    Task<bool> Create(Product product);
    Task<Product> GetEntityById(long id);
    Task<bool> Insert(Product product);
    Task<bool> SaveChanges();
}
