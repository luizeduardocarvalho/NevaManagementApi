namespace NevaManagement.Domain.Dtos.Organism;

public class EditOrganismDto
{
    [Required]
    public long? Id { get; set; }

    [Required]
    public string Name { get; set; }

    public string Description { get; set; }
}
