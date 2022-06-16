namespace NevaManagement.Domain.Models;

[Table("EquipmentUsage")]
public class EquipmentUsage : BaseEntity
{
    public Researcher Researcher { get; set; }

    public long ResearcherId { get; set; }

    public Equipment Equipment { get; set; }

    public long EquipmentId { get; set; }

    public DateTimeOffset StartDate { get; set; }

    public DateTimeOffset EndDate { get; set; }

    public string Description { get; set; }
}
