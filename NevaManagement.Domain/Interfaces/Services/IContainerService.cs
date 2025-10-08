namespace NevaManagement.Domain.Interfaces.Services;

public interface IContainerService
{
    Task<IList<GetSimpleContainerDto>> GetContainers(long laboratoryId);
    Task<bool> AddContainer(AddContainerDto addContainerDto);
    Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id, long laboratoryId);
    Task<GetDetailedContainerDto> GetDetailedContainer(long id, long laboratoryId);
    Task<IList<GetContainersByTransferDateDto>> GetContainersOrderedByTransferDate(int page, long laboratoryId);
}
