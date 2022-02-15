using Microsoft.EntityFrameworkCore;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Repositories
{
    public class ResearcherRepository : BaseRepository<Researcher>, IResearcherRepository
    {
        private readonly NevaManagementDbContext context;

        public ResearcherRepository(NevaManagementDbContext context)
            : base(context)
        {
            this.context = context;
        }

        public async Task<IEnumerable<Researcher>> GetAll()
        {
            return await this.context.Researchers.ToListAsync();
        }

        public async Task<bool> Create(Researcher researcher)
        {
            await Insert(researcher);
            return await SaveChanges();
        }
    }
}
