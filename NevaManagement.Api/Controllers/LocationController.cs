using NevaManagement.Domain.Dtos.Location;

namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class LocationController : ControllerBase
{
    private readonly ILocationService service;

    public LocationController(ILocationService service)
    {
        this.service = service;
    }

    [HttpGet("GetLocations")]
    public async Task<IActionResult> GetLocations()
    {
        var locations = await this.service.GetLocations();

        return Ok(locations);
    }

    [HttpPost("AddLocation")]
    public async Task<IActionResult> AddLocation([FromBody] AddLocationDto addLocationDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        var result = await this.service.AddLocation(addLocationDto);

        return result ?
            StatusCode(201, $"Successfully created {addLocationDto.Name}.") :
            StatusCode(500, "An error occurred while creating the location.");
    }

    [HttpPatch("EditLocation")]
    public async Task<IActionResult> EditProduct([FromBody] EditLocationDto editLocationDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        var result = await this.service.EditLocation(editLocationDto);

        return result ?
            Ok($"Successfully edited {editLocationDto.Name}.") :
            StatusCode(500, "An error occurred while editing the location.");
    }

    [HttpGet("GetLocationById")]
    public async Task<IActionResult> GetLocationById([FromQuery] long locationId)
    {
        if (locationId <= 0)
        {
            return BadRequest();
        }

        var location = await this.service.GetLocationById(locationId);

        return Ok(location);
    }
}
