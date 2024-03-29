﻿namespace NevaManagement.Api.Controllers;

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
    public async Task<IActionResult> AddOrganism([FromBody] AddOrganismDto addOrganismDto)
    {
        var result = await this.service.AddOrganism(addOrganismDto);

        return result ?
            StatusCode(201, $"Successfully created {addOrganismDto.Name}.") :
            StatusCode(500, "An error occurred while creating the organism.");
    }

    [HttpPatch("EditOrganism")]
    public async Task<IActionResult> EditOrganism([FromBody] EditOrganismDto editOrganismDto)
    {
        var result = await this.service.EditOrganism(editOrganismDto);

        return result ?
            Ok($"Successfully edited {editOrganismDto.Name}.") :
            StatusCode(500, "An error occurred while editing the organism.");
    }

    [HttpGet("GetOrganismById")]
    public async Task<IActionResult> GetOrganismById([FromQuery] long organismId)
    {
        var location = await this.service.GetOrganismById(organismId);

        return Ok(location);
    }
}
