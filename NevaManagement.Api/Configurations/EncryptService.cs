namespace NevaManagement.Api.Configurations;

public class EncryptService : IEncryptService
{
    private readonly IOptions<Settings> settings;

    public EncryptService(IOptions<Settings> settings)
    {
        this.settings = settings;
    }

    public string Encrypt(string password)
    {
        var key = this.settings.Value.Secret;
        if (string.IsNullOrEmpty(key))
            key = Environment.GetEnvironmentVariable("Settings");

        var keyBytes = Encoding.UTF8.GetBytes(key);

        var passwordBytes = Encoding.UTF8.GetBytes(password);
        var hash = new System.Security.Cryptography.HMACSHA256()
        {
            Key = keyBytes
        };

        var hashedPassword = hash.ComputeHash(passwordBytes);

        return Convert.ToBase64String(hashedPassword);
    }
}

