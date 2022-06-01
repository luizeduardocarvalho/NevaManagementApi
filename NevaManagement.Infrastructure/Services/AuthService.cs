namespace NevaManagement.Infrastructure.Services;

public class AuthService : IAuthService
{
    private readonly IResearcherRepository repository;

    public AuthService(IResearcherRepository repository)
    {
        this.repository = repository;
    }

    public async Task<bool> Register(RegisterDto registerDto)
    {
        var user = new Researcher()
        {
            Email = registerDto.Email,
            Name = registerDto.Name,
            Password = registerDto.Password,
            Role = registerDto.Role
        };

        try
        {
            return await this.repository.Create(user);
        }
        catch
        {
            throw new Exception("An error occurred while creating the new researcher.");
        }
    }
}