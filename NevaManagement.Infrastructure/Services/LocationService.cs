namespace NevaManagement.Infrastructure.Services;

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

        if (addLocationDto.SublocationId is not null)
        {
            var subLocation = await this.repository.GetById(addLocationDto.SublocationId.Value);
            location.SubLocation = subLocation;
        }

        try
        {
            await this.repository.Insert(location);
            result = await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while adding the location.");
        }

        return result;
    }

    public async Task<bool> EditLocation(EditLocationDto editLocationDto)
    {
        var result = false;
        var location = await this.repository.GetById(editLocationDto.Id);

        if (location is not null)
        {
            location.Name = editLocationDto.Name;
            location.Description = editLocationDto.Description;

            if (editLocationDto.SublocationId is not null)
            {
                var subLocation = await this.repository.GetById(editLocationDto.SublocationId.Value);
                location.SubLocation = subLocation;
            }

            try
            {
                result = await this.repository.SaveChanges();
            }
            catch
            {
                throw new Exception("An error occurred while updating the location.");
            }

        }

        return result;
    }

    public async Task<GetDetailedLocationDto> GetLocationById(long locationId)
    {
        try
        {
            return await this.repository.GetByDetailedLocationId(locationId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the location.");
        }
    }
}
