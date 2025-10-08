namespace NevaManagement.Infrastructure.Repositories;

public class ResearcherRepository : BaseRepository<Researcher>, IResearcherRepository
{
    private readonly NevaManagementDbContext context;

    public ResearcherRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IEnumerable<Researcher>> GetAll(long laboratoryId)
    {
        return await this.context.Researchers
            .Where(r => r.LaboratoryId == laboratoryId)
            .ToListAsync();
    }

    public async Task<IList<GetSimpleResearcherDto>> GetResearchers(long laboratoryId)
    {
        return await this.context.Researchers
                .Where(x => x.LaboratoryId == laboratoryId)
                .Select(x => new GetSimpleResearcherDto
                {
                    Id = x.Id,
                    Name = x.Name
                })
                .ToListAsync();
    }

    public async Task<GetDetailedResearcherDto> GetByEmailAndPassword(string email, string password)
    {
        return await this.context.Researchers
            .Where(researcher => researcher.Email.Equals(email) && researcher.Password.Equals(password))
            .Include(r => r.Laboratory)
            .Select(r => new GetDetailedResearcherDto()
            {
                Email = r.Email,
                Id = r.Id,
                Name = r.Name,
                Role = r.Role,
                LaboratoryId = r.LaboratoryId,
                LaboratoryName = r.Laboratory != null ? r.Laboratory.Name : null
            })
            .FirstOrDefaultAsync();
    }
}
