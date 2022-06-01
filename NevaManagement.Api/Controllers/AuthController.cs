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
    public async Task<ActionResult<dynamic>> AuthenticateAsync([FromBody] LoginDto loginDto)
    {
        var encryptedPassword = this.encryptService.Encrypt(loginDto.Password);
        var user = await this.researcherService.GetByEmailAndPassword(loginDto.Email, encryptedPassword);

        if (user == null)
            return NotFound("Email or password are invalid.");

        var token = this.tokenService.GenerateToken(user);

        user.Password = "";

        return Ok(new
        {
            user.Id,
            user.Name,
            user.Email,
            token,
            user
        });
    }

    [AllowAnonymous]
    [HttpPost("Register")]
    public async Task<IActionResult> Register([FromBody] RegisterDto registerDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        if (string.IsNullOrEmpty(registerDto.Password))
        {
            registerDto.Password = "123456";
        }

        registerDto.Password = this.encryptService.Encrypt(registerDto.Password);
        await this.authService.Register(registerDto);
        return Ok("Registered.");
    }
}
