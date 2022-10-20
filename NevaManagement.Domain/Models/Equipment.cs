namespace NevaManagement.Domain.Models;

[Table("Equipment")]
public class Equipment : BaseEntity
{
    public string Name { get; set; }

    public string Description { get; set; }

    [ForeignKey("Location_Id")]
    public Location Location { get; set; }

    public string PropertyNumber { get; set; }

    [ForeignKey("EquipmentId")]
    public IList<EquipmentUsage> EquipmentUsages { get; set; }
}
