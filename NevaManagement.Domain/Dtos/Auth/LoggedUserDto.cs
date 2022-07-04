namespace NevaManagement.Domain.Dtos.Auth;

public class LoggedUserDto
{
    public GetDetailedResearcherDto Researcher { get; set; }

    public string Token { get; set; }
}
