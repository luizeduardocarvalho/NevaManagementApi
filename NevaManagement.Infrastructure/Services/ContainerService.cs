﻿namespace NevaManagement.Infrastructure.Services;

public class ContainerService : IContainerService
{
    private readonly IContainerRepository repository;
    private readonly IResearcherRepository researcherRepository;
    private readonly IOrganismRepository organismRepository;

    public ContainerService(
        IContainerRepository repository,
        IResearcherRepository researcherRepository,
        IOrganismRepository organismRepository)
    {
        this.repository = repository;
        this.researcherRepository = researcherRepository;
        this.organismRepository = organismRepository;
    }

    public async Task<bool> AddContainer(AddContainerDto addContainerDto)
    {
        var result = false;
        var container = new Container
        {
            Description = addContainerDto.Description,
            Name = addContainerDto.Name,
            CreationDate = addContainerDto.Date,
            CultureMedia = addContainerDto.CultureMedia,
        };

        try
        {
            if (addContainerDto.SubContainerId is not null)
            {
                var subContainer = await this.repository.GetEntityById(addContainerDto.SubContainerId.Value);
                container.SubContainer = subContainer;
            }

            var researcher = await this.researcherRepository.GetEntityById(addContainerDto.ResearcherId.Value);
            if (researcher is not null)
            {
                container.Researcher = researcher;
            }

            var organism = await this.organismRepository.GetEntityById(addContainerDto.OrganismId.Value);
            if (organism is not null)
            {
                container.Origin = organism;
            }

            result = await this.repository.Create(container);
        }
        catch
        {
            throw;
        }

        return result;
    }

    public async Task<IList<GetSimpleContainerDto>> GetContainers()
    {
        return await this.repository.GetContainers();
    }

    public async Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id)
    {
        try
        {
            return await this.repository.GetChildrenContainers(id);
        }
        catch
        {
            throw;
        }
    }

    public async Task<GetDetailedContainerDto> GetDetailedContainer(long id)
    {
        try
        {
            return await this.repository.GetDetailedContainer(id);
        }
        catch
        {
            throw;
        }
    }
}
