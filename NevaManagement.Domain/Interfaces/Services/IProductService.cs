namespace NevaManagement.Domain.Interfaces.Services;

public interface IProductService
{
    Task<IEnumerable<GetProductDto>> GetAll();
    Task<GetProductDto> GetById(long id);
    Task<GetDetailedProductDto> GetDetailedProductById(long id);
    Task<bool> Create(CreateProductDto productDto);
    Task<bool> AddQuantityToProduct(AddQuantityToProductDto addQuantityToProductDto);
    Task<bool> UseProduct(UseProductDto useProductDto);
    Task<bool> EditProduct(EditProductDto editProductDto);
}
