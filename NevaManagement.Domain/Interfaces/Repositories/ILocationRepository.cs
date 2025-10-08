namespace NevaManagement.Domain.Interfaces.Repositories;

public interface ILocationRepository : IBaseRepository<Location>
{
    Task<GetDetailedLocationDto> GetByDetailedLocationId(long id, long laboratoryId);
    Task<IList<GetLocationDto>> GetLocations(long laboratoryId);
}
