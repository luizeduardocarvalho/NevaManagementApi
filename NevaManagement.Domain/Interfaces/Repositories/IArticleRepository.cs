namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IArticleRepository : IBaseRepository<Article>
{
    Task<GetArticleDto> GetArticle(long id);
    Task<IList<Article>> GetArticles(string[] dois);
}
