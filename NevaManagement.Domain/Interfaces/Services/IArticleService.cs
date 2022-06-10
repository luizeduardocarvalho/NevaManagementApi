namespace NevaManagement.Domain.Interfaces.Services;

public interface IArticleService
{
    Task<GetArticleDto> GetArticle(long id);
    Task<IList<Article>> GetArticles(string[] ids);
}
