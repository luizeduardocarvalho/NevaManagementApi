using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
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

        public async Task<GetDetailedLocationDto> GetById(long id)
        {
            return await this.context.Locations
                                        .Where(x => x.Id == id)
                                        .Include(x => x.SubLocation)
                                        .Select(x => new GetDetailedLocationDto
                                        {
                                            Name = x.Name,
                                            Description = x.Description,
                                            SubLocationId = x.SubLocation.Id,
                                            SubLocationName = x.SubLocation.Name
                                        })
                                        .FirstOrDefaultAsync();
        }

        public async Task<Location> GetEntityById(long id)
        {
            return await this.context.Locations
                                        .Where(x => x.Id == id)
                                        .FirstOrDefaultAsync();
        }

        public async Task<IList<GetLocationDto>> GetLocations()
        {
            return await this.context.Locations
                .Select(x => new GetLocationDto
                {
                    Id = x.Id,
                    Name = x.Name
                })
                .OrderBy(x => x.Id)
                .ToListAsync();
        }

        public async Task<bool> Create(Location location)
        {
            await Insert(location);
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
