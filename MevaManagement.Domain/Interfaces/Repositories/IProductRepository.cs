using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IProductRepository
    {
        Task<IEnumerable<GetProductDto>> GetAll();
        Task<GetDetailedProductDto> GetById(long id);
        Task<bool> Create(Product product);
    }
}
