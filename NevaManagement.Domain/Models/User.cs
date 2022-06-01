namespace NevaManagement.Domain.Models;

public class User : BaseEntity
{
    public string Email { get; set; }
    public string Password { get; set; }
    public string Role { get; set; }
}

