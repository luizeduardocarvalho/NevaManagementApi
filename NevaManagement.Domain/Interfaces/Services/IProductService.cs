namespace NevaManagement.Domain.Interfaces.Services;

public interface IProductService
{
    Task<IEnumerable<GetProductDto>> GetAll(int page, long laboratoryId);
    Task<GetProductDto> GetById(long id, long laboratoryId);
    Task<GetDetailedProductDto> GetDetailedProductById(long id, long laboratoryId);
    Task<bool> Create(CreateProductDto productDto);
    Task<bool> AddQuantityToProduct(AddQuantityToProductDto addQuantityToProductDto);
    Task<bool> UseProduct(UseProductDto useProductDto);
    Task<bool> EditProduct(EditProductDto editProductDto);
    Task<IList<GetDetailedProductDto>> GetLowInStockProducts();
}
