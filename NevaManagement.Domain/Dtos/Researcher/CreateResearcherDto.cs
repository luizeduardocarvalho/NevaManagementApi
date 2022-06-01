namespace NevaManagement.Domain.Dtos.Researcher;

public class CreateResearcherDto
{
    [Required]
    public string Name { get; set; }

    [Required]
    [EmailAddress]
    public string Email { get; set; }
}
