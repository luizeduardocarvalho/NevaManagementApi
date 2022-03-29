using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Interfaces.Services;
using System;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
    [ApiController]
    [Route("[controller]")]
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

            try
            {
                var result = await this.service.AddLocation(addLocationDto);

                if(result)
                {
                    return StatusCode(200, $"Successfully created {addLocationDto.Name}");
                }

                return StatusCode(500, "An error occurred while creating the location.");
            }
            catch(Exception ex)
            {
                return StatusCode(500, ex.Message);
            }
        }
    }
}
