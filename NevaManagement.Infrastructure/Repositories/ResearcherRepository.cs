namespace NevaManagement.Infrastructure.Repositories;

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

    public async Task<IList<GetSimpleResearcherDto>> GetResearchers()
    {
        return await this.context.Researchers
                .Select(x => new GetSimpleResearcherDto
                {
                    Id = x.Id,
                    Name = x.Name
                })
                .ToListAsync();
    }

    public async Task<Researcher> GetByEmailAndPassword(string email, string password)
    {
        return await this.context.Researchers
            .Where(researcher => researcher.Email.Equals(email) && researcher.Password.Equals(password))
            .FirstOrDefaultAsync();
    }
}
