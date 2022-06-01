namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IContainerRepository
{
    Task<IList<GetSimpleContainerDto>> GetContainers();
    Task<Container> GetEntityById(long id);
    Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id);
    Task<GetDetailedContainerDto> GetDetailedContainer(long id);
    Task<bool> Create(Container container);
    Task<bool> SaveChanges();
}
