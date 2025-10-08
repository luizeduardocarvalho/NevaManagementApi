namespace NevaManagement.Infrastructure.Services;

public class ContainerService : IContainerService
{
    private readonly IContainerRepository repository;
    private readonly IResearcherRepository researcherRepository;
    private readonly IOrganismRepository organismRepository;
    private readonly IArticleService articleService;
    private readonly IArticleContainerRepository articleContainerRepository;

    public ContainerService(
        IContainerRepository repository,
        IResearcherRepository researcherRepository,
        IOrganismRepository organismRepository,
        IArticleService articleService,
        IArticleContainerRepository articleContainerRepository)
    {
        this.repository = repository;
        this.researcherRepository = researcherRepository;
        this.organismRepository = organismRepository;
        this.articleService = articleService;
        this.articleContainerRepository = articleContainerRepository;
    }

    public async Task<bool> AddContainer(AddContainerDto addContainerDto)
    {
        var container = new Container
        {
            Description = addContainerDto.Description,
            Name = addContainerDto.Name,
            CreationDate = addContainerDto.CreationDate.UtcDateTime,
            TransferDate = addContainerDto.TransferDate.UtcDateTime,
            CultureMedia = addContainerDto.CultureMedia,
        };

        try
        {

            if (addContainerDto.SubContainerId is not null)
            {
                var subContainer = await this.repository.GetEntityById(addContainerDto.SubContainerId.Value, addContainerDto.LaboratoryId);
                if (subContainer is not null)
                {
                    container.SubContainer = subContainer;
                    container.ArticleContainerList = subContainer.ArticleContainerList;
                }
            }

            var researcher = await this.researcherRepository.GetById(addContainerDto.ResearcherId);
            if (researcher is not null)
            {
                container.Researcher = researcher;
            }

            var organism = await this.organismRepository.GetById(addContainerDto.OrganismId.Value);

            if (organism is not null)
            {
                container.Origin = organism;
            }

            if (addContainerDto.DoiList is not null)
            {
                var articles = await this.articleService.GetArticles(addContainerDto.DoiList.ToArray());

                if (!articles.Any())
                {
                    foreach (var doi in addContainerDto.DoiList)
                    {
                        var article = new Article()
                        {
                            Doi = doi
                        };

                        articles.Add(article);
                    }

                }

                foreach (var article in articles)
                {
                    var articleContainer = new ArticleContainer()
                    {
                        ArticleId = article.Id,
                        Container = container
                    };

                    await this.articleContainerRepository.Insert(articleContainer);
                }

                return await this.articleContainerRepository.SaveChanges();
            }
            else
            {
                await this.repository.Insert(container);
                return await this.repository.SaveChanges();
            }
        }
        catch
        {
            throw new Exception($"An error occurred while creating {addContainerDto.Name}.");
        }
    }

    public async Task<IList<GetSimpleContainerDto>> GetContainers(long laboratoryId)
    {
        try
        {
            return await this.repository.GetContainers(laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting all containers.");
        }
    }

    public async Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id, long laboratoryId)
    {
        try
        {
            return await this.repository.GetChildrenContainers(id, laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the children for the container.");
        }
    }

    public async Task<GetDetailedContainerDto> GetDetailedContainer(long id, long laboratoryId)
    {
        try
        {
            return await this.repository.GetDetailedContainer(id, laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the container.");
        }
    }

    public async Task<IList<GetContainersByTransferDateDto>> GetContainersOrderedByTransferDate(int page, long laboratoryId)
    {
        try
        {
            return await this.repository.GetContainersOrderedByTransferDate(page, laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the containers.");
        }
    }
}
