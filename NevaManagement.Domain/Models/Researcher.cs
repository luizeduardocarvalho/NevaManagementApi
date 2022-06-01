namespace NevaManagement.Domain.Models;

[Table("Researcher")]
public class Researcher : User
{
    [MaxLength(80)]
    public string Name { get; set; }
}
