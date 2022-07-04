namespace NevaManagement.Domain.Interfaces.Services;

public interface ITokenService
{
    Task<string> GenerateToken(GetDetailedResearcherDto user);
}
