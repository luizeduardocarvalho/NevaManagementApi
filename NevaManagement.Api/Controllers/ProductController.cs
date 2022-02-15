using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Interfaces.Services;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class ProductController : ControllerBase
    {
        private readonly IProductService service;

        public ProductController(IProductService service)
        {
            this.service = service;
        }

        [HttpGet("GetAll")]
        public async Task<IActionResult> GetAll()
        {
            var products = await this.service.GetAll();

            return Ok(products);
        }

        [HttpGet("GetById")]
        public async Task<IActionResult> GetById([FromQuery] long id)
        {
            var product = await this.service.GetById(id);

            return Ok(product);
        }

        [HttpPost("Create")]
        public async Task<IActionResult> Create([FromBody] CreateProductDto productDto)
        {
            var result = await this.service.Create(productDto);

            if(result)
            {
                return Ok("Success");
            }

            return StatusCode(500, "Error");
        }
    }
}
