using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.Container;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class ContainerRepository : BaseRepository<Container>, IContainerRepository
    {
        private readonly NevaManagementDbContext context;

        public ContainerRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<bool> Create(Container container)
        {
            await Insert(container);
            return await SaveChanges();
        }

        public async Task<IList<GetSimpleContainerDto>> GetContainers()
        {
            return await this.context.Containers
                .Select(x => new GetSimpleContainerDto 
                { 
                    Id = x.Id,
                    Name = x.Name
                })
                .ToListAsync();
        }

        public async Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id)
        {
            return await this.context.Containers
                .Where(container => container.SubContainer.Id == id)
                .Select(container => new GetSimpleContainerDto 
                { 
                    Id = container.Id,
                    Name = container.Name
                })
                .ToListAsync();
        }

        public async Task<GetDetailedContainerDto> GetDetailedContainer(long id)
        {
            return await this.context.Containers
                .Where(container => container.Id == id)
                .Include(container => container.Researcher)
                .Include(container => container.Origin)
                .Select(container => new GetDetailedContainerDto
                {
                    CreationDate = container.CreationDate,
                    CultureMedia = container.CultureMedia,
                    Description = container.Description,
                    Name = container.Name,
                    OriginName = container.Origin.Name,
                    ResearcherName = container.Researcher.Name
                })
                .FirstOrDefaultAsync();
        }

        public async Task<Container> GetEntityById(long id)
        {
            return await this.context.Containers.Where(x => x.Id == id).SingleOrDefaultAsync();
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
