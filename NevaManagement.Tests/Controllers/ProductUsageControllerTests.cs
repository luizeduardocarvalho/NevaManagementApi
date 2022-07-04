using NevaManagement.Domain.Dtos.ProductUsage;

namespace NevaManagement.Tests.Controllers;

public sealed class ProductUsageControllerTests
{
    private readonly Mock<IProductUsageService> productUsageServiceMock;
    private readonly ProductUsageController controller;

    public ProductUsageControllerTests()
    {
        productUsageServiceMock = new Mock<IProductUsageService>();

        controller = new ProductUsageController(productUsageServiceMock.Object);
    }

    [Fact]
    public async Task GetLastUsesByResearcherShouldReturnOk()
    {
        // Arrange
        const int id = 1;
        var lastUses = new List<GetLastUsesByResearcherDto>()
        {
            new()
            {
                Name = "LastUser",
                Quantity = 1,
                Unit = "ml"
            }
        };

        productUsageServiceMock.Setup(s => s.GetLastUsesByResearcher(id)).ReturnsAsync(lastUses);

        // Act
        var result = await this.controller.GetLastUsesByResearcher(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productUsageServiceMock.Verify(s => s.GetLastUsedProductByResearcher(id), Times.Once);
    }

    [Fact]
    public async Task GetLastUsedProductByResearcherShouldReturnOk()
    {
        // Arrange
        const int id = 1;
        var lastUses = new GetLastUsedProductDto()
        {
            Quantity = 1,
            Unit = "ml"
        };

        productUsageServiceMock
            .Setup(s => s.GetLastUsedProductByResearcher(id))
            .ReturnsAsync(lastUses);

        // Act
        var result = await this.controller.GetLastUsedProductByResearcher(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productUsageServiceMock.Verify(s => s.GetLastUsedProductByResearcher(id), Times.Once);
    }

    [Fact]
    public async Task GetLastUsesByProductShouldReturnOk()
    {
        // Arrange
        const int id = 1;
        var lastUses = new List<GetLastUseByProductDto>()
        {
            new()
            {
                Quantity = 1,
                Unit = "ml"
            }
        };

        productUsageServiceMock.Setup(s => s.GetLastUsesByProduct(id)).ReturnsAsync(lastUses);

        // Act
        var result = await this.controller.GetLastUsesByProduct(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productUsageServiceMock.Verify(s => s.GetLastUsesByProduct(id), Times.Once);
    }
}
