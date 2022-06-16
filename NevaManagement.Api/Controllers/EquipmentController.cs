namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class EquipmentController : ControllerBase
{
    private readonly IEquipmentService service;

    public EquipmentController(IEquipmentService service)
    {
        this.service = service;
    }

    [HttpGet("GetEquipments")]
    public async Task<IActionResult> GetEquipments()
    {
        var result = await this.service.GetEquipments();

        return Ok(result);
    }

    [HttpPost("AddEquipment")]
    public async Task<IActionResult> AddEquipment([FromBody] AddEquipmentDto equipmentDto)
    {
        var result = await this.service.AddEquipment(equipmentDto);

        return Ok(result);
    }

    [HttpGet("GetDetailedEquipment")]
    public async Task<IActionResult> GetDetailedEquipment([FromQuery] long id)
    {
        var result = await this.service.GetDetailedEquipment(id);

        return Ok("The equipment was updated successfully.");
    }

    [HttpPatch("EditEquipment")]
    public async Task<IActionResult> EditEquipment([FromBody] EditEquipmentDto equipmentDto)
    {
        var result = await this.service.EditEquipment(equipmentDto);

        if(result)
        {
            return Ok($"{equipmentDto.Name} was updated successfully.");
        }

        return StatusCode(500, $"An error occurred while updating {equipmentDto.Name}.");
    }
}
