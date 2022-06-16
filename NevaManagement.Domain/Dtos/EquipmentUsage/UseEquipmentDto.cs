namespace NevaManagement.Domain.Dtos.EquipmentUsage;

public class UseEquipmentDto
{
    [Required]
    public long EquipmentId { get; set; }

    [Required]
    public long ResearcherId { get; set; }

    public string Description { get; set; }

    [Required]
    public DateTimeOffset StartDate { get; set; }

    [Required]
    public DateTimeOffset EndDate { get; set; }
}
