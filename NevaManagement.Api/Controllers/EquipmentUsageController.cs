namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class EquipmentUsageController : ControllerBase
{
    private readonly IEquipmentUsageService service;

    public EquipmentUsageController(IEquipmentUsageService service)
    {
        this.service = service;
    }

    [HttpGet("GetEquipmentUsageCalendar")]
    public async Task<IActionResult> GetEquipmentUsageCalendar([FromQuery] long id, [FromQuery] long laboratoryId)
    {
        var result = await this.service.GetEquipmentUsageCalendar(id, laboratoryId);

        return Ok(result);
    }

    [HttpGet("GetEquipmentUsageHistory")]
    public async Task<IActionResult> GetEquipmentUsageHistory([FromQuery] long id, [FromQuery] long laboratoryId)
    {
        var result = await this.service.GetEquipmentUsageHistory(id, laboratoryId);

        return Ok(result);
    }

    [HttpPost("UseEquipment")]
    public async Task<IActionResult> UseEquipment(UseEquipmentDto useEquipmentDto)
    {
        var result = await this.service.UseEquipment(useEquipmentDto);

        return result ?
            Ok("Equipment was used successfully.") :
            StatusCode(500, "An error occurred while using the equipment.");
    }

    [HttpGet("GetEquipmentUsage")]
    public async Task<IActionResult> GetEquipmentUsage(long equipmentId, int page, long laboratoryId)
    {
        var result = await this.service.GetEquipmentUsage(equipmentId, page, laboratoryId);

        return Ok(result);
    }
}
