using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Dtos.Product;
using NevaManagement.Domain.Interfaces.Services;
using System;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
    [ApiController]
    [Route("[controller]")]
    [Produces("application/json")]
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

        [HttpGet("GetDetailedProductById")]
        public async Task<IActionResult> GetDetailedProductById([FromQuery] long id)
        {
            var product = await this.service.GetDetailedProductById(id);

            return Ok(product);
        }

        [HttpGet("GetProductById")]
        public async Task<IActionResult> GetProductById([FromQuery] long id)
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

        [HttpPatch("AddQuantity")]
        public async Task<IActionResult> AddProductQuantity([FromBody] AddQuantityToProductDto addQuantityToProductDto)
        {
            if (!ModelState.IsValid)
            {
                return BadRequest();
            }

            try
            {
                var result = await this.service.AddQuantityToProduct(addQuantityToProductDto);

                if(result)
                {
                    return StatusCode(200, $"Successfully added {addQuantityToProductDto.Quantity} to product.");
                }

                return StatusCode(500, "An error occurred while adding quantity to product.");
            }
            catch(Exception ex)
            {
                return StatusCode(500, ex.Message);
            }
        }

        [HttpPatch("UseProduct")]
        public async Task<IActionResult> UseProduct([FromBody] UseProductDto useProductDto)
        {
            if(!ModelState.IsValid)
            {
                return BadRequest();
            }

            try
            {
                var result = await this.service.UseProduct(useProductDto);

                if(result)
                {
                    return Ok($"Successfully used {useProductDto.Quantity}.");
                }
            }
            catch(Exception ex)
            {
                return StatusCode(500, ex.Message);
            }

            return StatusCode(500, "Error occurred while using product.");
        }

        [HttpPatch("EditProduct")]
        public async Task<IActionResult> EditProduct([FromBody] EditProductDto editProductDto)
        {
            if (!ModelState.IsValid)
            {
                return BadRequest();
            }

            try
            {
                var result = await this.service.EditProduct(editProductDto);

                if (result)
                {
                    return StatusCode(200, $"Successfully edited {editProductDto.Name}.");
                }
            }
            catch (Exception ex)
            {
                return StatusCode(500, ex.Message);
            }

            return StatusCode(500, "Error occurred while editing the product.");
        }
    }
}
