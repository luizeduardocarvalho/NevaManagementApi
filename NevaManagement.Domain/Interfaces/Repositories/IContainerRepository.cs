namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IContainerRepository : IBaseRepository<Container>
{
    Task<IList<GetSimpleContainerDto>> GetContainers();
    Task<Container> GetEntityById(long id);
    Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id);
    Task<GetDetailedContainerDto> GetDetailedContainer(long id);
    Task<IList<GetContainersByTransferDateDto>> GetContainersOrderedByTransferDate(int page);
}
