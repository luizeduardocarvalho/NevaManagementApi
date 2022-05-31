namespace NevaManagement.Domain.Models;

[Table("Article")]
public class Article : BaseEntity
{
    public string Doi { get; set; }

    [ForeignKey("ArticleId")]
    public IList<ArticleContainer> ArticleContainerList { get; set; }
}