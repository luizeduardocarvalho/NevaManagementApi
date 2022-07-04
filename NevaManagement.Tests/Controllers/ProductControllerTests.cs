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

    [Fact]
    public async Task GetProductByIdShouldReturnOk()
    {
        // Arrange
        var id = 1;
        var expectedResult = new GetProductDto
        {
            Id = id,
            Name = "Name"
        };

        productServiceMock
            .Setup(service => service.GetById(id))
            .ReturnsAsync(expectedResult);

        // Act
        var result = await this.controller.GetProductById(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.GetById(id), Times.Once);
    }

    [Fact]
    public async Task CreateProductShouldReturnOk()
    {
        // Arrange
        var product = new CreateProductDto()
        {
            Name = "Name"
        };

        productServiceMock
            .Setup(service => service.Create(product))
            .ReturnsAsync(true);

        // Act
        var result = await this.controller.Create(product);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.Create(product), Times.Once);
    }

    [Fact]
    public async Task CreateProductWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var product = new CreateProductDto()
        {
            Name = "Name"
        };

        productServiceMock
            .Setup(service => service.Create(product))
            .ReturnsAsync(false);

        // Act
        var result = await this.controller.Create(product);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        productServiceMock
            .Verify(service => service.Create(product), Times.Once);
    }

    [Fact]
    public async Task AddQuantityShouldReturnOk()
    {
        // Arrange
        var quantity = new AddQuantityToProductDto()
        {
            ProductId = 1,
            Quantity = 1
        };

        productServiceMock
            .Setup(service => service.AddQuantityToProduct(quantity))
            .ReturnsAsync(true);

        // Act
        var result = await this.controller.AddProductQuantity(quantity);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.AddQuantityToProduct(quantity), Times.Once);
    }

    [Fact]
    public async Task AddQuantityWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var quantity = new AddQuantityToProductDto()
        {
            ProductId = 1,
            Quantity = 1
        };

        productServiceMock
            .Setup(service => service.AddQuantityToProduct(quantity))
            .ReturnsAsync(false);

        // Act
        var result = await this.controller.AddProductQuantity(quantity);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        productServiceMock
            .Verify(service => service.AddQuantityToProduct(quantity), Times.Once);
    }

    [Fact]
    public async Task UseProductShouldReturnOk()
    {
        // Arrange
        var useProduct = new UseProductDto()
        {
            ProductId = 1,
            Quantity = 1
        };

        productServiceMock
            .Setup(service => service.UseProduct(useProduct))
            .ReturnsAsync(true);

        // Act
        var result = await this.controller.UseProduct(useProduct);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.UseProduct(useProduct), Times.Once);
    }

    [Fact]
    public async Task UseProductWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var useProduct = new UseProductDto()
        {
            ProductId = 1,
            Quantity = 1
        };

        productServiceMock
            .Setup(service => service.UseProduct(useProduct))
            .ReturnsAsync(false);

        // Act
        var result = await this.controller.UseProduct(useProduct);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        productServiceMock
            .Verify(service => service.UseProduct(useProduct), Times.Once);
    }

    [Fact]
    public async Task EditProductShouldReturnOk()
    {
        // Arrange
        var product = new EditProductDto()
        {
            Id = 1,
            Name = "Product"
        };

        productServiceMock
            .Setup(service => service.EditProduct(product))
            .ReturnsAsync(true);

        // Act
        var result = await this.controller.EditProduct(product);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        productServiceMock
            .Verify(service => service.EditProduct(product), Times.Once);
    }

    [Fact]
    public async Task EditProductWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var product = new EditProductDto()
        {
            Id = 1,
            Name = "Product"
        };

        productServiceMock
            .Setup(service => service.EditProduct(product))
            .ReturnsAsync(false);

        // Act
        var result = await this.controller.EditProduct(product);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        productServiceMock
            .Verify(service => service.EditProduct(product), Times.Once);
    }
}
