using NevaManagement.Domain.Dtos.Location;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Services
{
    public interface ILocationService
    {
        Task<IList<GetLocationDto>> GetLocations();
        Task<bool> AddLocation(AddLocationDto addLocationDto);
    }
}
