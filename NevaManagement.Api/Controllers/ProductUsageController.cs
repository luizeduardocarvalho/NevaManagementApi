namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
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
        if (researcherId == 0)
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

    [HttpGet("GetLastUsesByProduct")]
    public async Task<IActionResult> GetLastUsesByProduct([FromQuery] long productId)
    {
        if (productId == 0)
        {
            return BadRequest();
        }

        var result = await this.service.GetLastUsesByProduct(productId);

        return Ok(result);
    }
}
