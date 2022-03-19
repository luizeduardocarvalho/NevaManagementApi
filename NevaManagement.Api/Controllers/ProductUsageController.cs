using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Interfaces.Services;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
    [ApiController]
    public class ProductUsageController : ControllerBase
    {
        private readonly IProductUsageService service;

        public ProductUsageController(IProductUsageService service)
        {
            this.service = service;
        }

        [HttpGet("GetLastUsesByResearcher")]
        public async Task<IActionResult> GetLastUsesByResearcher([FromQuery] long researcherId)
        {
            if(researcherId == 0)
            {
                return BadRequest();
            }

            var result = await this.service.GetLastUsesByResearcher(researcherId);

            return Ok(result);
        }

        [HttpGet("GetLastUsedProductByResearcher")]
        public async Task<IActionResult> GetLastUsedProductByResearcher([FromQuery] long researcherId)
        {
            if (researcherId == 0)
            {
                return BadRequest();
            }

            var result = await this.service.GetLastUsedProductByResearcher(researcherId);

            return Ok(result);
        }
    }
}
