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

    public async Task<bool> Create(Researcher researcher)
    {
        await Insert(researcher);
        return await SaveChanges();
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

    public async Task<Researcher> GetEntityById(long id)
    {
        return await this.context.Researchers
            .Where(x => x.Id == id)
            .SingleOrDefaultAsync();
    }

    public new async Task<bool> SaveChanges()
    {
        var result = await this.context.SaveChangesAsync();

        if (result > 0)
        {
            return true;
        }

        return false;
    }
}
