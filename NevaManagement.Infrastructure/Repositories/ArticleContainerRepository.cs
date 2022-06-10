namespace NevaManagement.Infrastructure.Repositories;

public class ArticleContainerRepository : BaseRepository<ArticleContainer>, IArticleContainerRepository
{
    public ArticleContainerRepository(NevaManagementDbContext context)
        : base(context)
    {
    }
}
