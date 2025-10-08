namespace NevaManagement.Infrastructure.Repositories;

public class ContainerRepository : BaseRepository<Container>, IContainerRepository
{
    private readonly NevaManagementDbContext context;

    public ContainerRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<IList<GetSimpleContainerDto>> GetContainers(long laboratoryId)
    {
        return await this.context.Containers
            .Where(x => x.LaboratoryId == laboratoryId)
            .Select(x => new GetSimpleContainerDto
            {
                Id = x.Id,
                Name = x.Name
            })
            .ToListAsync();
    }

    public async Task<IList<GetSimpleContainerDto>> GetChildrenContainers(long id, long laboratoryId)
    {
        return await this.context.Containers
            .Where(container => container.SubContainer.Id == id && container.LaboratoryId == laboratoryId)
            .Select(container => new GetSimpleContainerDto
            {
                Id = container.Id,
                Name = container.Name
            })
            .ToListAsync();
    }

    public async Task<GetDetailedContainerDto> GetDetailedContainer(long id, long laboratoryId)
    {
        return await this.context.Containers
            .Where(container => container.Id == id && container.LaboratoryId == laboratoryId)
            .Include(container => container.ArticleContainerList)
            .Include(container => container.Researcher)
            .Include(container => container.Origin)
            .Select(container => new GetDetailedContainerDto
            {
                CreationDate = container.CreationDate,
                CultureMedia = container.CultureMedia,
                Description = container.Description,
                Name = container.Name,
                OriginName = container.Origin.Name,
                ResearcherName = container.Researcher.Name,
                DoiList = container.ArticleContainerList
                    .Select(doi => doi.Article.Doi)
                    .ToList()
            })
            .FirstOrDefaultAsync();
    }

    public async Task<Container> GetEntityById(long id, long laboratoryId)
    {
        return await this.context.Containers
            .Include(x => x.ArticleContainerList)
            .ThenInclude(x => x.Article)
            .Where(x => x.Id == id && x.LaboratoryId == laboratoryId).SingleOrDefaultAsync();
    }

    public async Task<IList<GetContainersByTransferDateDto>> GetContainersOrderedByTransferDate(int page, long laboratoryId)
    {
        return await this.context.Containers
            .Where(x => x.LaboratoryId == laboratoryId)
            .OrderBy(x => x.Name)
            .Skip((page - 1) * 15)
            .Take(15)
            .Where(x => x.SubContainer == null)
            .OrderBy(x => x.TransferDate)
            .Select(x => new GetContainersByTransferDateDto()
            {
                Id = x.Id,
                TransferDate = x.TransferDate,
                Name = x.Name
            })
            .ToListAsync();
    }
}
