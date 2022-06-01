namespace NevaManagement.Domain.Interfaces.Services;

public interface ITokenService
{
    string GenerateToken(User user);
}
