namespace NevaManagement.Domain.Interfaces.Services;

public interface IEncryptService
{
    Task<string> Encrypt(string password);
}
