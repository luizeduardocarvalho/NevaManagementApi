namespace NevaManagement.Domain.Dtos.Equipment;

public class EditEquipmentDto
{
    [Required]
    public long? Id { get; set; }

    public string Description { get; set; }

    public string Name { get; set; }

    public string Patrimony { get; set; }

    public long LocationId { get; set; }
}

