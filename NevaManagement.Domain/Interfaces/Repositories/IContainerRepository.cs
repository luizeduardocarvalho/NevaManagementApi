namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IContainerRepository : IBaseRepository<Container>
{
    Task<IList<GetSimpleContainerDto>> GetContainers(long laboratoryId);
    Task<Container> GetEntityById(long id, long laboratoryId);
    Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id, long laboratoryId);
    Task<GetDetailedContainerDto> GetDetailedContainer(long id, long laboratoryId);
    Task<IList<GetContainersByTransferDateDto>> GetContainersOrderedByTransferDate(int page, long laboratoryId);
}
