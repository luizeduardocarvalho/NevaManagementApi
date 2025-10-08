namespace NevaManagement.Domain.Models;

[Table("Researcher")]
public class Researcher : User
{
    [MaxLength(80)]
    public string Name { get; set; }

    [ForeignKey("ResearcherId")]
    public IList<EquipmentUsage> EquipmentUsages { get; set; }

    [ForeignKey("Laboratory_Id")]
    public Laboratory Laboratory { get; set; }
    public long LaboratoryId { get; set; }
}
