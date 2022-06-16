namespace NevaManagement.Domain.Dtos.Equipment;

public class AddEquipmentDto
{
    [Required]
    public string Name { get; set; }

    public string Description { get; set; }
}
