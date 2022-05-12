using NevaManagement.Domain.Dtos.Organism;

namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class OrganismController : ControllerBase
{
    private readonly IOrganismService service;

    public OrganismController(IOrganismService service)
    {
        this.service = service;
    }

    [HttpGet("GetOrganisms")]
    public async Task<IActionResult> GetOrganisms()
    {
        var organisms = await this.service.GetOrganisms();

        return Ok(organisms);
    }

    [HttpPost("AddOrganism")]
    public async Task<IActionResult> AddLocation([FromBody] AddOrganismDto addOrganismDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        try
        {
            var result = await this.service.AddOrganism(addOrganismDto);

            if (result)
            {
                return StatusCode(201, $"Successfully created {addOrganismDto.Name}.");
            }

            return StatusCode(500, "An error occurred while creating the organism.");
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    [HttpPatch("EditOrganism")]
    public async Task<IActionResult> EditProduct([FromBody] EditOrganismDto editOrganismDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        try
        {
            var result = await this.service.EditOrganism(editOrganismDto);

            if (result)
            {
                return StatusCode(200, $"Successfully edited {editOrganismDto.Name}.");
            }

            return StatusCode(500, "An error occurred while editing the organism.");
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    [HttpGet("GetOrganismById")]
    public async Task<IActionResult> GetLocationById([FromQuery] long organismId)
    {
        if (organismId <= 0)
        {
            return BadRequest();
        }

        var location = await this.service.GetOrganismById(organismId);

        if (location is not null)
        {
            return Ok(location);
        }

        return StatusCode(500, $"Error finding organism with id {organismId}");
    }
}
