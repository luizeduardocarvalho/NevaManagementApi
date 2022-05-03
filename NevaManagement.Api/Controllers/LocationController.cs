using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Dtos.Location;
using NevaManagement.Domain.Interfaces.Services;
using System;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
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

            try
            {
                var result = await this.service.AddLocation(addLocationDto);

                if(result)
                {
                    return StatusCode(201, $"Successfully created {addLocationDto.Name}.");
                }

                return StatusCode(500, "An error occurred while creating the location.");
            }
            catch(Exception ex)
            {
                return StatusCode(500, ex.Message);
            }
        }

        [HttpPatch("EditLocation")]
        public async Task<IActionResult> EditProduct([FromBody] EditLocationDto editLocationDto)
        {
            if (!ModelState.IsValid)
            {
                return BadRequest();
            }

            try
            {
                var result = await this.service.EditLocation(editLocationDto);

                if(result)
                {
                    return StatusCode(200, $"Successfully edited {editLocationDto.Name}.");
                }

                return StatusCode(500, "An error occurred while editing the location.");
            }
            catch(Exception ex)
            {
                return StatusCode(500, ex.Message);
            }
        }

        [HttpGet("GetLocationById")]
        public async Task<IActionResult> GetLocationById([FromQuery] long locationId)
        {
            if(locationId <= 0)
            {
                return BadRequest();
            }

            var location = await this.service.GetLocationById(locationId);

            if(location is not null)
            {
                return Ok(location);
            }

            return StatusCode(500, $"Error finding location with id {locationId}");
        }
    }
}
