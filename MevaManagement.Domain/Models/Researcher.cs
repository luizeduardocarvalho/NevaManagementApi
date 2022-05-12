namespace NevaManagement.Domain.Models;

[Table("Researcher")]
public class Researcher : BaseEntity
{
    [MaxLength(80)]
    public string Name { get; set; }

    [MaxLength(100)]
    public string Email { get; set; }
}
