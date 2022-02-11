using NevaManagement.Domain.Dtos.Researcher;
using NevaManagement.Domain.Models;
using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Services
{
    public interface IResearcherService
    {
        Task<Researcher> GetAll();
        Task<bool> Create(CreateResearcherDto researcherDto);
    }
}
