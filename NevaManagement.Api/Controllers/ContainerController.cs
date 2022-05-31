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
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetSimpleContainerDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetContainers()
    {
        var locations = await this.service.GetContainers();
        return Ok(locations);
    }

    [HttpPost("AddContainer")]
    [ProducesResponseType(StatusCodes.Status201Created, Type = typeof(string))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> AddContainer([FromBody] AddContainerDto addContainerDto)
    {
        if (!ModelState.IsValid)
        {
            return BadRequest();
        }

        await this.service.AddContainer(addContainerDto);

        return StatusCode(201, $"Successfully created {addContainerDto.Name}.");
    }

    [HttpGet("GetChildrenContainers")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetSimpleContainerDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetChildrenContainers([FromQuery] long containerId)
    {
        if (containerId <= 0)
        {
            return BadRequest();
        }

        var result = await this.service.GetChildrenContainers(containerId);

        return Ok(result);
    }

    [HttpGet("GetDetailedContainer")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetDetailedContainerDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetDetailedContainer([FromQuery] long containerId)
    {
        if (containerId <= 0)
        {
            return BadRequest();
        }

        var result = await this.service.GetDetailedContainer(containerId);

        return Ok(result);
    }
}
