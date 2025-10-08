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
    public async Task<IActionResult> GetContainers([FromQuery] long laboratoryId)
    {
        var locations = await this.service.GetContainers(laboratoryId);
        return Ok(locations);
    }

    [HttpPost("AddContainer")]
    [ProducesResponseType(StatusCodes.Status201Created, Type = typeof(string))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> AddContainer([FromBody] AddContainerDto addContainerDto)
    {
        await this.service.AddContainer(addContainerDto);

        return StatusCode(201, $"Successfully created {addContainerDto.Name}.");
    }

    [HttpGet("GetChildrenContainers")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetSimpleContainerDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetChildrenContainers([FromQuery] long containerId, [FromQuery] long laboratoryId)
    {
        var result = await this.service.GetChildrenContainers(containerId, laboratoryId);

        return Ok(result);
    }

    [HttpGet("GetDetailedContainer")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetDetailedContainerDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetDetailedContainer([FromQuery] long containerId, [FromQuery] long laboratoryId)
    {
        var result = await this.service.GetDetailedContainer(containerId, laboratoryId);

        return Ok(result);
    }

    [HttpGet("GetContainersOrderedByTransferDate")]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(IList<GetContainersByTransferDateDto>))]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status500InternalServerError)]
    public async Task<IActionResult> GetContainersOrderedByTransferDate([FromQuery] int page, [FromQuery] long laboratoryId)
    {
        var result = await this.service.GetContainersOrderedByTransferDate(page, laboratoryId);

        return Ok(result);
    }
}
