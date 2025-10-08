namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IProductRepository : IBaseRepository<Product>
{
    Task<IEnumerable<GetProductDto>> GetAll(int page, long laboratoryId);
    Task<GetProductDto> GetProductById(long id, long laboratoryId);
    Task<GetDetailedProductDto> GetDetailedProductById(long id, long laboratoryId);
}
