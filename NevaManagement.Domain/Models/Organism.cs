namespace NevaManagement.Domain.Models;

[Table("Organism")]
public class Organism : BaseEntity
{
    [MaxLength(100)]
    public string Name { get; set; }

    public string Type { get; set; }

    public string Description { get; set; }

    [Column("Collection_Date")]
    public DateTimeOffset CollectionDate { get; set; }

    [Column("Collection_Location")]
    public string CollectionLocation { get; set; }

    [Column("Isolation_Date")]
    public DateTimeOffset IsolationDate { get; set; }

    [ForeignKey("Origin_Id")]
    public Organism Origin { get; set; }

    [Column("Origin_Part")]
    public string OriginPart { get; set; }

    [ForeignKey("Laboratory_Id")]
    public Laboratory Laboratory { get; set; }
    public long LaboratoryId { get; set; }
}
