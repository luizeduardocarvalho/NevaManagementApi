namespace NevaManagement.Domain.Dtos.Auth;

public class RegisterDto
{
    [Required]
    public string Email { get; set; }

    [Required]
    public string Name { get; set; }

    [Required]
    public string Password { get; set; }

    [Required]
    public string Role { get; set; }
}
