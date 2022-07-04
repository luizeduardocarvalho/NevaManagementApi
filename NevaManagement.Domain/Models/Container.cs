namespace NevaManagement.Domain.Models;

[Table("Container")]
public class Container : BaseEntity
{
    [Column("Creation_Date")]
    public DateTimeOffset CreationDate { get; set; }

    [Column("Culture_Media")]
    public string CultureMedia { get; set; }

    public string Description { get; set; }

    [ForeignKey("ContainerId")]
    public IList<ArticleContainer> ArticleContainerList { get; set; }

    [ForeignKey("Equipment_Id")]
    public Equipment Equipment { get; set; }

    [MaxLength(80)]
    public string Name { get; set; }

    [ForeignKey("Origin_Id")]
    public Organism Origin { get; set; }

    [ForeignKey("Researcher_Id")]
    public Researcher Researcher { get; set; }

    [ForeignKey("Sub_Container_Id")]
    public Container SubContainer { get; set; }

    [Column("TransferDate")]
    public DateTimeOffset TransferDate { get; set; }
}
