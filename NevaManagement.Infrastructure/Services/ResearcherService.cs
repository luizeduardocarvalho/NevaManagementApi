namespace NevaManagement.Infrastructure.Services;

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

    public async Task<IList<GetSimpleResearcherDto>> GetResearchers()
    {
        return await this.repository.GetResearchers();
    }

    public async Task<Researcher> GetByEmailAndPassword(string email, string password)
    {
        try
        {
            return await repository.GetByEmailAndPassword(email, password);
        }
        catch
        {
            throw new Exception("User not found");
        }
    }
}
