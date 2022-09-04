namespace NevaManagement.Infrastructure.Services;

public class OrganismService : IOrganismService
{
    private readonly IOrganismRepository repository;
    private readonly IBuilder<AddOrganismDto, Organism> builder;

    public OrganismService(
        IOrganismRepository repository,
        IBuilder<AddOrganismDto, Organism> builder)
    {
        this.repository = repository;
        this.builder = builder;
    }

    public async Task<IList<GetOrganismDto>> GetOrganisms()
    {
        return await this.repository.GetOrganisms();
    }

    public async Task<bool> AddOrganism(AddOrganismDto addOrganismDto)
    {
        var organism = this.builder.Build(addOrganismDto);

        if (addOrganismDto.OriginOrganismId is not null)
        {
            var originOrganism = await this.repository.GetById(addOrganismDto.OriginOrganismId.Value);
            organism.Origin = originOrganism;
        }

        try
        {
            await this.repository.Insert(organism);
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while adding the organism.");
        }
    }

    public async Task<bool> EditOrganism(EditOrganismDto editOrganismDto)
    {
        var result = false;
        var organism = await this.repository.GetById(editOrganismDto.Id.Value);

        if (organism is not null)
        {
            organism.Name = editOrganismDto.Name;
            organism.Description = editOrganismDto.Description;

            try
            {
                result = await this.repository.SaveChanges();
            }
            catch
            {
                throw;
            }

        }

        return result;
    }

    public async Task<GetDetailedOrganismDto> GetOrganismById(long organismId)
    {
        try
        {
            return await this.repository.GetDetailedOrganismById(organismId);
        }
        catch
        {
            throw new Exception("An error occurred while getting organism.");
        }
    }
}
