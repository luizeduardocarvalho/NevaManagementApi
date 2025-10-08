using Microsoft.Extensions.Logging;

namespace NevaManagement.Api.Configurations;

public class TokenService : ITokenService
{
    private readonly IOptions<Settings> settings;
    private readonly ILogger<TokenService> logger;

    public TokenService(IOptions<Settings> settings, ILogger<TokenService> logger)
    {
        this.settings = settings;
        this.logger = logger;
    }

    public async Task<string> GenerateToken(GetDetailedResearcherDto user)
    {
        try
        {
            var secret = Environment.GetEnvironmentVariable("Settings") 
                ?? this.settings.Value.Secret;

            var tokenHandler = new JwtSecurityTokenHandler();
            var key = Encoding.ASCII.GetBytes(secret);
            var tokenDescriptor = new SecurityTokenDescriptor
            {
                Subject = new ClaimsIdentity(new[]
                {
                        new Claim(ClaimTypes.Name, user.Email),
                        new Claim(ClaimTypes.Role, user.Role),
                        new Claim(ClaimTypes.NameIdentifier, user.Id.ToString()),
                        new Claim("LaboratoryId", user.LaboratoryId?.ToString() ?? ""),
                        new Claim("LaboratoryName", user.LaboratoryName ?? "")
                    }),
                Expires = DateTime.UtcNow.AddDays(30),
                SigningCredentials = new SigningCredentials(
                    new SymmetricSecurityKey(key),
                    SecurityAlgorithms.HmacSha256Signature)
            };

            var token = tokenHandler.CreateToken(tokenDescriptor);

            return await Task.FromResult(tokenHandler.WriteToken(token));
        }
        catch(Exception e)
        {
            logger.LogError(e, e.Message);
            throw;
        }
    }
}
