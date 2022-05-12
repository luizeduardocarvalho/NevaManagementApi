namespace NevaManagement.Domain.Models;

[Table("ProductUsage")]
public class ProductUsage : BaseEntity
{
    [ForeignKey("Researcher_Id")]
    public Researcher Researcher { get; set; }

    [ForeignKey("Product_Id")]
    public Product Product { get; set; }

    [Column("Usage_Date")]
    public DateTimeOffset UsageDate { get; set; }

    public double Quantity { get; set; }

    public string Description { get; set; }
}
