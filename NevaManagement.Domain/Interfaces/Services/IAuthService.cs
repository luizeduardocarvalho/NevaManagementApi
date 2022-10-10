namespace NevaManagement.Domain.Interfaces.Services;

public interface IAuthService
{
    Task ChangePassword(ChangePasswordDto changePasswordDto);
}
