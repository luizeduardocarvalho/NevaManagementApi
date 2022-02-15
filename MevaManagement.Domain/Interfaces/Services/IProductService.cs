using NevaManagement.Domain.Dtos.Product;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Services
{
    public interface IProductService
    {
        Task<IEnumerable<GetProductDto>> GetAll();
        Task<GetDetailedProductDto> GetById(long id);
        Task<bool> Create(CreateProductDto productDto);
    }
}
