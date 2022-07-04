namespace NevaManagement.Api.Configurations;

public class TokenService : ITokenService
{
    private readonly IOptions<Settings> settings;

    public TokenService(IOptions<Settings> settings)
    {
        this.settings = settings;
    }

    public async Task<string> GenerateToken(GetDetailedResearcherDto user)
    {
        var secret = this.settings.Value.Secret;
        if (string.IsNullOrEmpty(secret))
            secret = Environment.GetEnvironmentVariable("Settings");

        var tokenHandler = new JwtSecurityTokenHandler();
        var key = Encoding.ASCII.GetBytes(secret);
        var tokenDescriptor = new SecurityTokenDescriptor
        {
            Subject = new ClaimsIdentity(new[]
            {
                    new Claim(ClaimTypes.Name, user.Email),
                    new Claim(ClaimTypes.Role, user.Role)
                }),
            Expires = DateTime.UtcNow.AddDays(30),
            SigningCredentials = new SigningCredentials(
                new SymmetricSecurityKey(key),
                SecurityAlgorithms.HmacSha256Signature)
        };

        var token = tokenHandler.CreateToken(tokenDescriptor);

        return await Task.FromResult(tokenHandler.WriteToken(token));
    }
}
