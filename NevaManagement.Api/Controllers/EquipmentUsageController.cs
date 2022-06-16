using NevaManagement.Domain.Dtos.EquipmentUsage;
using NevaManagement.Domain.Dtos.Product;

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
    public async Task<IActionResult> GetEquipmentUsageCalendar([FromQuery] long id)
    {
        var result = await this.service.GetEquipmentUsageCalendar(id);

        return Ok(result);
    }

    [HttpGet("GetEquipmentUsageHistory")]
    public async Task<IActionResult> GetEquipmentUsageHistory([FromQuery] long id)
    {
        var result = await this.service.GetEquipmentUsageHistory(id);

        return Ok(result);
    }

    [HttpPost("UseEquipment")]
    public async Task<IActionResult> UseEquipment(UseEquipmentDto useEquipmentDto)
    {
        var result = await this.service.UseEquipment(useEquipmentDto);

        return Ok(result);
    }
}
