namespace NevaManagement.Infrastructure.Repositories;

public class ArticleRepository : BaseRepository<Article>, IArticleRepository
{
    private readonly NevaManagementDbContext context;

    public ArticleRepository(NevaManagementDbContext context)
        : base(context)
    {
        this.context = context;
    }

    public async Task<GetArticleDto> GetArticle(long id)
    {
        return await this.context.Articles
            .Where(article => article.Id == id)
            .Select(article => new GetArticleDto()
            {
                Id = article.Id,
                Doi = article.Doi
            })
            .FirstOrDefaultAsync();
    }

    public async Task<IList<Article>> GetArticles(string[] dois)
    {
        return await this.context.Articles
            .Where(article => dois.Contains(article.Doi))
            .ToListAsync();
    }
}
