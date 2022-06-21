namespace NevaManagement.Infrastructure.Services;

public class ResearcherService : IResearcherService
{
    private readonly IResearcherRepository repository;

    public ResearcherService(IResearcherRepository repository)
    {
        this.repository = repository;
    }

    public async Task<IList<GetSimpleResearcherDto>> GetResearchers()
    {
        try
        {
            return await this.repository.GetResearchers();
        }
        catch
        {
            throw new Exception("An error occurred while getting researchers.");
        }
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

    public async Task<bool> Create(CreateResearcherDto researcher)
    {
        var newResearcher = new Researcher()
        {
            Name = researcher.Name,
            Email = researcher.Email
        };

        try
        {
            await this.repository.Insert(newResearcher);
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while creating the researcher.");
        }
    }
}
