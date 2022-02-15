using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
    public class ProductService : IProductService
    {
        private readonly IProductRepository repository;
        private readonly ILocationRepository locationRepository;

        public ProductService(IProductRepository repository, ILocationRepository locationRepository)
        {
            this.repository = repository;
            this.locationRepository = locationRepository;
        }

        public async Task<IEnumerable<GetProductDto>> GetAll()
        {
            return await this.repository.GetAll();
        }

        public async Task<GetDetailedProductDto> GetById(long id)
        {
            return await this.repository.GetById(id);
        }

        public async Task<bool> Create(CreateProductDto productDto)
        {
            var location =  await this.locationRepository.GetEntityById(productDto.LocationId);
            var product = new Product
            {
                Description = productDto.Description,
                Location = location,
                Name = productDto.Name,
                Quantity = productDto.Quantity,
                Unit = productDto.Unit
            };

            return await this.repository.Create(product);
        }
    }
}
