using Microsoft.AspNetCore.Mvc.ModelBinding;

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

        return result ?
            Ok($"{equipmentDto.Name} was created successfully.") :
            StatusCode(500, $"An error occurred while creating {equipmentDto.Name}.");
    }

    [HttpGet("GetDetailedEquipment")]
    public async Task<IActionResult> GetDetailedEquipment([BindRequired, FromQuery] long id)
    {
        var result = await this.service.GetDetailedEquipment(id);

        return Ok(result);
    }

    [HttpPatch("EditEquipment")]
    public async Task<IActionResult> EditEquipment([FromBody] EditEquipmentDto equipmentDto)
    {
        var result = await this.service.EditEquipment(equipmentDto);

        return result ?
            Ok($"{equipmentDto.Name} was updated successfully.") :
            StatusCode(500, $"An error occurred while updating {equipmentDto.Name}.");
    }
}
