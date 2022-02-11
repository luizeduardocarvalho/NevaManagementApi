using NevaManagement.Domain.Dtos.Researcher;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Domain.Models;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
    public class ResearcherService : IResearcherService
    {
        private readonly IResearcherRepository repository;

        public ResearcherService(IResearcherRepository repository)
        {
            this.repository = repository;
        }

        public Task<bool> Create(CreateResearcherDto researcherDto)
        {
            throw new System.NotImplementedException();
        }

        public Task<Researcher> GetAll()
        {
            throw new System.NotImplementedException();
        }
    }
}
