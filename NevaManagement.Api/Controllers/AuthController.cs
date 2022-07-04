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
    [HttpPost("Register")]
    public async Task<IActionResult> Register([FromBody] RegisterDto registerDto)
    {
        if (string.IsNullOrEmpty(registerDto.Password))
        {
            registerDto.Password = "123456";
        }

        registerDto.Password = await this.encryptService.Encrypt(registerDto.Password);
        var result = await this.authService.Register(registerDto);

        return result ?
            Ok("The user was registered successfully.") :
            StatusCode(500, "An error occurred while registering the user.");
    }
}
