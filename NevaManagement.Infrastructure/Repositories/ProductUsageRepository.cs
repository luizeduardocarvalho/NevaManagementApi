using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.ProductUsage;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class ProductUsageRepository : BaseRepository<ProductUsage>, IProductUsageRepository
    {
        private readonly NevaManagementDbContext context;

        public ProductUsageRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<IList<GetLastUsesByResearcherDto>> GetLastUsesByResearcher(long researcherId)
        {
            return await this.context.ProductUsages
                .Where(x => x.Researcher.Id == researcherId)
                .Select(x => new GetLastUsesByResearcherDto
                {
                    Quantity = x.Quantity,
                    ResearcherName = x.Researcher.Name,
                    Unit = x.Product.Unit,
                    UsageDate = x.UsageDate
                })
                .ToListAsync();
        }

        public async Task<IList<GetLastUseByProductDto>> GetLastUsesByProduct(long productId)
        {
            return await this.context.ProductUsages
                .Where(x => x.Product.Id == productId)
                .Select(x => new GetLastUseByProductDto
                {
                    Quantity = x.Quantity,
                    ResearcherName = x.Researcher.Name,
                    Unit = x.Product.Unit,
                    UseDate = x.UsageDate
                })
                .ToListAsync();
        }

        public async Task<GetLastUsedProductDto> GetLastUsedProductByResearcher(long researcherId)
        {
            return await this.context.ProductUsages
                .Where(x => x.Id == researcherId)
                .OrderBy(x => x.UsageDate)
                .Select(x => new GetLastUsedProductDto
                {
                    ProductId = x.Product.Id,
                    LocationName = x.Product.Location.Name,
                    ProductName = x.Product.Name,
                    Quantity = x.Quantity,
                    Unit = x.Product.Unit
                })
                .LastOrDefaultAsync();
        }

        public new async Task<bool> SaveChanges()
        {
            var result = await this.context.SaveChangesAsync();

            if (result > 0)
            {
                return true;
            }

            return false;
        }

        public async Task<bool> Create(ProductUsage productUsage)
        {
            await Insert(productUsage);
            return await SaveChanges();
        }
    }
}
