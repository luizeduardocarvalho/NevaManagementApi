namespace NevaManagement.Domain.Interfaces.Repositories;

public interface ILocationRepository : IBaseRepository<Location>
{
    Task<GetDetailedLocationDto> GetByDetailedLocationId(long id);
    Task<IList<GetLocationDto>> GetLocations();
}
