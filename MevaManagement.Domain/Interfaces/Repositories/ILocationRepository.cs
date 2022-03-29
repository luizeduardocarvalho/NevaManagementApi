using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface ILocationRepository
    {
        Task<GetLocationDto> GetById(long id);
        Task<Location> GetEntityById(long id);
        Task<IList<GetLocationDto>> GetLocations();
        Task<bool> Create(Location location);
    }
}
