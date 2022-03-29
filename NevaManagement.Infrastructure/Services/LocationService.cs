using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
    public class LocationService : ILocationService
    {
        private readonly ILocationRepository repository;

        public LocationService(ILocationRepository repository)
        {
            this.repository = repository;
        }

        public async Task<IList<GetLocationDto>> GetLocations()
        {
            return await this.repository.GetLocations();
        }

        public async Task<bool> AddLocation(AddLocationDto addLocationDto)
        {
            var result = false;
            var location = new Location
            {
                Description = addLocationDto.Description,
                Name = addLocationDto.Name
            };

            if(addLocationDto.SublocationId is not null)
            {
                var sublocation = await this.repository.GetEntityById(addLocationDto.SublocationId.Value);
                location.SubLocation = sublocation;
            }

            try
            {
                result = await this.repository.Create(location);
            }
            catch
            {
                throw;
            }

            return result;
        }
    }
}
