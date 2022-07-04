namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
public class ResearcherController : ControllerBase
{
    private readonly IResearcherService service;

    public ResearcherController(IResearcherService service)
    {
        this.service = service;
    }

    [HttpGet("GetResearchers")]
    public async Task<IActionResult> GetResearchers()
    {
        var researchers = await this.service.GetResearchers();

        return Ok(researchers);
    }

    [HttpPost]
    public async Task<IActionResult> Create(CreateResearcherDto researcherDto)
    {
        var result = await this.service.Create(researcherDto);

        return result ?
            Ok("Successfully create a new researcher.") :
            StatusCode(500, "Could not save the new researcher.");
    }
}
