namespace NevaManagement.Domain.Interfaces.Services;

public interface IOrganismService
{
    Task<IList<GetOrganismDto>> GetOrganisms(long laboratoryId);
    Task<bool> AddOrganism(AddOrganismDto addOrganismDto);
    Task<bool> EditOrganism(EditOrganismDto editOrganismDto);
    Task<GetDetailedOrganismDto> GetOrganismById(long organismId, long laboratoryId);
}
