using NevaManagement.Domain.Dtos.Product;

namespace NevaManagement.Tests.Controllers;

public sealed class ProductControllerTests
{
    private readonly Mock<IProductService> productServiceMock;
    private readonly ProductController controller;

    public ProductControllerTests()
    {
        productServiceMock = new Mock<IProductService>();
        controller = new ProductController(productServiceMock.Object);
    }

    [Fact]
    public async Task GetAllShouldReturnOk()
    {
        // Arrange
        var expectedResult = new List<GetProductDto>
        {
            new()
            {
                Id = 1,
                Name = "Name"
            }
        };

        productServiceMock
            .Setup(service => service.GetAll())
            .ReturnsAsync(expectedResult);

        // Act
        var result = await this.controller.GetAll();

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.GetAll(), Times.Once);
    }

    [Fact]
    public async Task GetDetailedProductByIdShouldReturnOk()
    {
        // Arrange
        var expectedResult = new GetDetailedProductDto()
        {
            Id = 1,
            Name = "Product"
        };

        productServiceMock
            .Setup(service => service.GetDetailedProductById(1))
            .ReturnsAsync(expectedResult);

        // Act
        var result = await this.controller.GetDetailedProductById(1);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.GetDetailedProductById(1), Times.Once);
    }
}