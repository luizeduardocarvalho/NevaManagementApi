namespace NevaManagement.Domain.Interfaces.Services;

public interface IAuthService
{
    Task<bool> Register(RegisterDto registerDto);
}
