namespace NevaManagement.Domain.Models;

[Table("Equipment")]
public class Equipment : BaseEntity
{
    public string Name { get; set; }

    public string Description { get; set; }

    [ForeignKey("Location_Id")]
    public Location Location { get; set; }

    public string PropertyNumber { get; set; }

    [ForeignKey("Laboratory_Id")]
    public Laboratory Laboratory { get; set; }
    public long LaboratoryId { get; set; }

    [ForeignKey("EquipmentId")]
    public IList<EquipmentUsage> EquipmentUsages { get; set; }
}
