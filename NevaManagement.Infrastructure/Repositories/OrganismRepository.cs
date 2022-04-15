using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.Organism;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class OrganismRepository : BaseRepository<Organism>, IOrganismRepository
    {
        private readonly NevaManagementDbContext context;

        public OrganismRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<GetDetailedOrganismDto> GetById(long id)
        {
            return await this.context.Organisms
                                        .Where(x => x.Id == id)
                                        .Include(x => x.Origin)
                                        .Select(x => new GetDetailedOrganismDto
                                        {
                                            Name = x.Name,
                                            Description = x.Description,
                                            CollectionDate = x.CollectionDate,
                                            CollectionLocation = x.CollectionLocation,
                                            IsolationDate = x.IsolationDate,
                                            OriginOrganismId = x.Origin.Id,
                                            OriginPart = x.OriginPart
                                        })
                                        .FirstOrDefaultAsync();
        }

        public async Task<Organism> GetEntityById(long id)
        {
            return await this.context.Organisms
                                        .Where(x => x.Id == id)
                                        .FirstOrDefaultAsync();
        }

        public async Task<IList<GetOrganismDto>> GetOrganisms()
        {
            return await this.context.Organisms
                .Select(x => new GetOrganismDto
                {
                    Id = x.Id,
                    Name = x.Name
                })
                .ToListAsync();
        }

        public async Task<bool> Create(Organism organism)
        {
            await Insert(organism);
            return await SaveChanges();
        }

        public async new Task<bool> SaveChanges()
        {
            var result = await this.context.SaveChangesAsync();

            if (result > 0)
            {
                return true;
            }

            return false;
        }
    }
}
