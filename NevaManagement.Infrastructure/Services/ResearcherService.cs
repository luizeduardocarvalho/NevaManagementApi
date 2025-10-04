namespace NevaManagement.Infrastructure.Services;

public class ResearcherService : IResearcherService
{
    private readonly IResearcherRepository repository;

    public ResearcherService(IResearcherRepository repository)
    {
        this.repository = repository;
    }

    async Task<IList<GetSimpleResearcherDto>> IResearcherService.GetResearchers()
    {
        try
        {
            return await this.repository.GetResearchers().ConfigureAwait(false);
        }
        catch
        {
            throw new Exception("An error occurred while getting researchers.");
        }
    }

    async Task<GetDetailedResearcherDto> IResearcherService.GetByEmailAndPassword(string email, string password)
    {
        try
        {
            return await repository
                .GetByEmailAndPassword(email, password)
                .ConfigureAwait(false);
        }
        catch
        {
            throw new Exception("User not found");
        }
    }

    async Task<bool> IResearcherService.Create(CreateResearcherDto researcher)
    {
        var newResearcher = new Researcher()
        {
            Name = researcher.Name,
            Email = researcher.Email,
            Password = researcher.Password,
            Role = researcher.Role ?? "Researcher"
        };

        try
        {
            await this.repository.Insert(newResearcher).ConfigureAwait(false);
            return await this.repository.SaveChanges().ConfigureAwait(false);
        }
        catch
        {
            throw new Exception("An error occurred while creating the researcher.");
        }
    }

    async Task IResearcherService.Save()
    {
        await this.repository.SaveChanges().ConfigureAwait(false);
    }

    async Task<Researcher> IResearcherService.GetById(long id)
    {
        return await this.repository.GetById(id);
    }
}
