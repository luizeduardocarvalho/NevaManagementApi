using NevaManagement.Domain.Dtos.Container;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Services
{
    public interface IContainerService
    {
        Task<IList<GetSimpleContainerDto>> GetContainers();
        Task<bool> AddContainer(AddContainerDto addContainerDto);
    }
}
