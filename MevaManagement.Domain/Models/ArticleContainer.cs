namespace NevaManagement.Domain.Models
{
    [Table("ArticleContainer")]
    public class ArticleContainer : BaseEntity
    {
        public long ArticleId { get; set; }

        public Article Article { get; set; }

        public long ContainerId { get; set; }

        public Container Container { get; set; }
    }
}
