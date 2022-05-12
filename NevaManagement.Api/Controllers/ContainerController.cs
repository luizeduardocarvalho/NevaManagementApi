namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class ContainerController : ControllerBase
{
    private readonly IContainerService service;

    public ContainerController(IContainerService service)
    {
        this.service = service;
    }

    [HttpGet("GetContainers")]
    public async Task<IActionResult> GetContainers()
    {
        var locations = await this.service.GetContainers();

        return Ok(locations);
    }

    [HttpPost("AddContainer")]
    public async Task<IActionResult> AddContainer([FromBody] AddContainerDto addContainerDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        try
        {
            var result = await this.service.AddContainer(addContainerDto);

            if (result)
            {
                return StatusCode(201, $"Successfully created {addContainerDto.Name}.");
            }

            return StatusCode(500, "An error occurred while creating the container.");
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    [HttpGet("GetChildrenContainers")]
    public async Task<IActionResult> GetChildrenCotainers([FromQuery] long containerId)
    {
        if (containerId <= 0)
        {
            return BadRequest();
        }

        try
        {
            var result = await this.service.GetChildrenContainers(containerId);

            return Ok(result);
        }
        catch(Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    [HttpGet("GetDetailedContainer")]
    public async Task<IActionResult> GetDetailedContainer([FromQuery] long containerId)
    {
        if (containerId <= 0)
        {
            return BadRequest();
        }

        try
        {
            var result = await this.service.GetDetailedContainer(containerId);

            return Ok(result);
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }
}
