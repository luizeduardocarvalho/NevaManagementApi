namespace NevaManagement.Api.Controllers;

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
    public async Task<IActionResult> GetAll([FromQuery] int page, [FromQuery] long laboratoryId)
    {
        var products = await this.service.GetAll(page, laboratoryId);

        return Ok(products);
    }

    [HttpGet("GetDetailedProductById")]
    public async Task<IActionResult> GetDetailedProductById([FromQuery] long id, [FromQuery] long laboratoryId)
    {
        var product = await this.service.GetDetailedProductById(id, laboratoryId);

        return Ok(product);
    }

    [HttpGet("GetProductById")]
    public async Task<IActionResult> GetProductById([FromQuery] long id, [FromQuery] long laboratoryId)
    {
        var product = await this.service.GetById(id, laboratoryId);

        return Ok(product);
    }

    [HttpPost("Create")]
    public async Task<IActionResult> Create([FromBody] CreateProductDto productDto)
    {
        var result = await this.service.Create(productDto);

        return result ? 
            StatusCode(201, $"Successfully created {productDto.Name}.") :
            StatusCode(500, "An error occurred while creating the product.");
    }

    [HttpPatch("AddQuantity")]
    public async Task<IActionResult> AddProductQuantity([FromBody] AddQuantityToProductDto addQuantityToProductDto)
    {
        var result = await this.service.AddQuantityToProduct(addQuantityToProductDto);

        return result ?
            StatusCode(200, $"Successfully added {addQuantityToProductDto.Quantity} to product.") :
            StatusCode(500, "An error occurred while adding quantity to product.");
    }

    [HttpPatch("UseProduct")]
    public async Task<IActionResult> UseProduct([FromBody] UseProductDto useProductDto)
    {
        var result = await this.service.UseProduct(useProductDto);

        return result ?
            Ok($"Successfully used {useProductDto.Quantity}{useProductDto.Unit}.") :
            StatusCode(500, "Error occurred while using product.");
    }

    [HttpPatch("EditProduct")]
    public async Task<IActionResult> EditProduct([FromBody] EditProductDto editProductDto)
    {
        var result = await this.service.EditProduct(editProductDto);

        return result ?
            StatusCode(200, $"Successfully edited {editProductDto.Name}.") :
            StatusCode(500, "Error occurred while editing the product.");
    }

    [HttpGet("GetLowInStockProducts")]
    public async Task<IActionResult> GetLowInStockProducts()
    {
        var result = await this.service.GetLowInStockProducts();

        return Ok(result);
    }
}
