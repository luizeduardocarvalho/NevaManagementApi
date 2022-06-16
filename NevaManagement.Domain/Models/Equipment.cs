namespace NevaManagement.Domain.Models;

[Table("Equipment")]
public class Equipment : BaseEntity
{
    public string Name { get; set; }

    public string Description { get; set; }

    [ForeignKey("EquipmentId")]
    public IList<EquipmentUsage> EquipmentUsages { get; set; }
}
