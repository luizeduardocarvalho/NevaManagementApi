namespace NevaManagement.Infrastructure.Services;

public class ArticleService : IArticleService
{
    private readonly IArticleRepository articleRepository;

    public ArticleService(IArticleRepository articleRepository)
    {
        this.articleRepository = articleRepository;
    }

    public async Task<GetArticleDto> GetArticle(long id)
    {
        return await this.articleRepository.GetArticle(id);
    }

    public async Task<IList<Article>> GetArticles(string[] ids)
    {
        return await this.articleRepository.GetArticles(ids);
    }
}
