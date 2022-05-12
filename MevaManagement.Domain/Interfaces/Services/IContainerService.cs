namespace NevaManagement.Domain.Interfaces.Services;

public interface IContainerService
{
    Task<IList<GetSimpleContainerDto>> GetContainers();
    Task<bool> AddContainer(AddContainerDto addContainerDto);
    Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id);
    Task<GetDetailedContainerDto> GetDetailedContainer(long id);
}
