namespace NevaManagement.Domain.Dtos.Container;

public class AddContainerDto
{
    public string[] DoiList { get; set; }
    public string CultureMedia { get; set; }
    [Required]
    public string Name { get; set; }

    public string Description { get; set; }

    [Required]
    public long? ResearcherId { get; set; }

    public long? SubContainerId { get; set; }

    [Required]
    public long? OrganismId { get; set; }


    public DateTimeOffset Date { get; set; }
}
