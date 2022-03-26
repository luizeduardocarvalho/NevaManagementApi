using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Domain.Models;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
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

            if(product is not null)
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
    }
}
