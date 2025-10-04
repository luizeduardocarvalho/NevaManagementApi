namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class AuthController : ControllerBase
{
    private readonly IEncryptService encryptService;
    private readonly IResearcherService researcherService;
    private readonly ITokenService tokenService;
    private readonly IAuthService authService;

    public AuthController(
        IEncryptService encryptService,
        IResearcherService researcherService,
        ITokenService tokenService,
        IAuthService authService)
    {
        this.encryptService = encryptService;
        this.researcherService = researcherService;
        this.tokenService = tokenService;
        this.authService = authService;
    }

    [AllowAnonymous]
    [HttpPost("login")]
    public async Task<IActionResult> AuthenticateAsync([FromBody] LoginDto loginDto)
    {
        var encryptedPassword = await this.encryptService.Encrypt(loginDto.Password);
        var user = await this.researcherService.GetByEmailAndPassword(loginDto.Email, encryptedPassword);

        if (user == null)
            return NotFound("Email or password are invalid.");

        var token = await this.tokenService.GenerateToken(user);

        user.Password = "";

        var loggedUser = new LoggedUserDto { Researcher = user, Token = token };

        return Ok(loggedUser);
    }

    [AllowAnonymous]
    [HttpPost("register")]
    public async Task<IActionResult> RegisterAsync([FromBody] RegisterDto registerDto)
    {
        var encryptedPassword = await this.encryptService.Encrypt(registerDto.Password);

        var researcher = new CreateResearcherDto
        {
            Name = registerDto.Name,
            Email = registerDto.Email,
            Password = encryptedPassword,
            Role = registerDto.Role
        };

        var created = await this.researcherService.Create(researcher);

        if (!created)
            return BadRequest("Failed to create user.");

        var user = await this.researcherService.GetByEmailAndPassword(registerDto.Email, encryptedPassword);
        var token = await this.tokenService.GenerateToken(user);

        user.Password = "";

        var loggedUser = new LoggedUserDto { Researcher = user, Token = token };

        return Ok(loggedUser);
    }

    [HttpPost("ChangePassword")]
    public async Task<IActionResult> ChangePassword([FromBody] ChangePasswordDto changePasswordDto)
    {
        await this.authService.ChangePassword(changePasswordDto);

        return NoContent();
    }
}
