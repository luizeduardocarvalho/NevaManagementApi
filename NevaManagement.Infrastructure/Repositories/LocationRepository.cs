using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Linq;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class LocationRepository : BaseRepository<Location>, ILocationRepository
    {
        private readonly NevaManagementDbContext context;

        public LocationRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<GetLocationDto> GetById(long id)
        {
            return await this.context.Locations
                                        .Where(x => x.Id == id)
                                        .Select(x => new GetLocationDto
                                        {
                                            Id = x.Id,
                                            Name = x.Name
                                        })
                                        .FirstOrDefaultAsync();
        }

        public async Task<Location> GetEntityById(long id)
        {
            return await this.context.Locations
                                        .Where(x => x.Id == id)
                                        .FirstOrDefaultAsync();
        }
    }
}
