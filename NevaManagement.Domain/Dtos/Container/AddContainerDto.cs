namespace NevaManagement.Domain.Dtos.Container;

public class AddContainerDto
{
    public string[] DoiList { get; set; }

    public string CultureMedia { get; set; }

    [Required]
    public string Name { get; set; }

    public string Description { get; set; }

    [Required]
    public long ResearcherId { get; set; }

    public long? SubContainerId { get; set; }

    public long? OrganismId { get; set; }

    [Required]
    public DateTimeOffset CreationDate { get; set; }

    [Required]
    public DateTimeOffset TransferDate { get; set; }

    [Required]
    public long LaboratoryId { get; set; }
}
