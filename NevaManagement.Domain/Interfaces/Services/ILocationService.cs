namespace NevaManagement.Domain.Interfaces.Services;

public interface ILocationService
{
    Task<IList<GetLocationDto>> GetCachedLocations();
    Task<bool> AddLocation(AddLocationDto addLocationDto);
    Task<bool> EditLocation(EditLocationDto editLocationDto);
    Task<GetDetailedLocationDto> GetLocationById(long locationId);
}
