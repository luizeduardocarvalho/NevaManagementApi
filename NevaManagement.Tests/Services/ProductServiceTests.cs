using NevaManagement.Domain.Dtos.Product;

namespace NevaManagement.Tests.Services;

public sealed class ProductServiceTests
{
    private readonly Mock<IProductRepository> productRepositoryMock;
    private readonly Mock<ILocationRepository> locationRepositoryMock;
    private readonly Mock<IResearcherRepository> researcherRepositoryMock;
    private readonly Mock<IProductUsageRepository> productUsageRepositoryMock;
    private ProductService service;

    public ProductServiceTests()
    {
        productRepositoryMock = new Mock<IProductRepository>();
        locationRepositoryMock = new Mock<ILocationRepository>();
        researcherRepositoryMock = new Mock<IResearcherRepository>();
        productUsageRepositoryMock = new Mock<IProductUsageRepository>();

        service = new ProductService(
            productRepositoryMock.Object,
            locationRepositoryMock.Object,
            researcherRepositoryMock.Object,
            productUsageRepositoryMock.Object);
    }

    [Fact]
    public async Task GetAll_ShouldCallGetAll()
    {
        // Arrange
        var expectedProducts = new List<GetProductDto>()
        {
            new()
            {
                Id = 1,
                Name = "Name"
            }
        };

        productRepositoryMock
            .Setup(repository => repository.GetAll())
            .ReturnsAsync(expectedProducts);

        // Act
        await service.GetAll();

        // Assert
        productRepositoryMock.Verify(repository => repository.GetAll(), Times.Once);
    }

    [Fact]
    public async Task GetAllWithDatabaseExceptionShouldReturnException()
    {
        // Arrange
        productRepositoryMock
            .Setup(repository => repository.GetAll())
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => service.GetAll();

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting all products.");
        productRepositoryMock.Verify(repository => repository.GetAll(), Times.Once);
    }

    [Theory]
    [InlineData(1)]
    [InlineData(20)]
    public async Task GetDetailedProductByIdShouldCallGetAll(long id)
    {
        // Arrange
        var expectedDetailedProducts= new GetDetailedProductDto()
        {
            Id = 1,
            Name = "Name"
        };

        productRepositoryMock
            .Setup(repository => repository.GetDetailedProductById(id))
            .ReturnsAsync(expectedDetailedProducts);

        // Act
        await service.GetDetailedProductById(id);

        // Assert
        productRepositoryMock.Verify(repository => repository.GetDetailedProductById(id), Times.Once);
    }

    [Theory]
    [InlineData(0)]
    public async Task GetDetailedProductByIdWithDatabaseExceptionShouldReturnException(long id)
    {
        // Arrange
        productRepositoryMock
            .Setup(repository => repository.GetDetailedProductById(id))
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => service.GetDetailedProductById(id);

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting a detailed product.");
        productRepositoryMock.Verify(repository => repository.GetDetailedProductById(id), Times.Once);
    }
}
