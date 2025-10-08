namespace NevaManagement.Infrastructure.Repositories;

public class LocationRepository : BaseRepository<Location>, ILocationRepository
{
    private readonly NevaManagementDbContext context;

    public LocationRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<GetDetailedLocationDto> GetByDetailedLocationId(long id, long laboratoryId)
    {
        return await this.context.Locations
            .Where(x => x.Id == id && x.LaboratoryId == laboratoryId)
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

    public async Task<IList<GetLocationDto>> GetLocations(long laboratoryId)
    {
        return await this.context.Locations
            .Where(x => x.LaboratoryId == laboratoryId)
            .Select(x => new GetLocationDto
            {
                Id = x.Id,
                Name = x.Name
            })
            .OrderBy(x => x.Id)
            .ToListAsync();
    }
}
