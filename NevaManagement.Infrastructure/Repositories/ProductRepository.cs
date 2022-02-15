using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class ProductRepository : BaseRepository<Product>, IProductRepository
    {
        private readonly NevaManagementDbContext context;

        public ProductRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<IEnumerable<GetProductDto>> GetAll()
        {
            return await this.context.Products
                                        .Select(p => new GetProductDto
                                        {
                                            Id = p.Id,
                                            Name = p.Name
                                        }).ToListAsync();
        }

        public async Task<GetDetailedProductDto> GetById(long id)
        {
            return await this.context.Products
                                        .Where(x => x.Id == id)
                                        .Include(x => x.Location)
                                        .Select(x => new GetDetailedProductDto
                                        {
                                            Id = x.Id,
                                            Location = new GetLocationDto
                                            {
                                                Id = x.Location.Id,
                                                Name = x.Location.Name
                                            },
                                            Name = x.Name,
                                            Quantity = x.Quantity,
                                            Unit = x.Unit
                                        })
                                        .FirstOrDefaultAsync();
        }

        public async Task<bool> Create(Product product)
        {
            await Insert(product);
            return await SaveChanges();
        }
    }
}
