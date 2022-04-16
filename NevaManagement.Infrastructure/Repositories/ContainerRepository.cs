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
