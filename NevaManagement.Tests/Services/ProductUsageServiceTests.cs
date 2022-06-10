using NevaManagement.Domain.Dtos.ProductUsage;

namespace NevaManagement.Tests.Services;

public sealed class ProductUsageServiceTests
{
    private readonly Mock<IProductUsageRepository> productRepositoryMock;
    private readonly ProductUsageService controller;

    public ProductUsageServiceTests()
    {
        productRepositoryMock = new Mock<IProductUsageRepository>();

        controller = new ProductUsageService(productRepositoryMock.Object);
    }

    [Fact]
    public async Task GetLastUsesByResearcherShouldCallRepository()
    {
        // Arrange
        var researcherId = 1;
        var getLastUsesByResearcherDto = new List<GetLastUsesByResearcherDto>()
        {
            new()
            {
                Name = "Name",
                Quantity = 1,
                Unit = "Unit",
                UsageDate = DateTimeOffset.Now
            }
        };

        productRepositoryMock
            .Setup(repository => repository.GetLastUsesByResearcher(researcherId))
            .ReturnsAsync(getLastUsesByResearcherDto);

        // Act
        var result = await this.controller.GetLastUsesByResearcher(researcherId);

        // Assert
        result.Should().BeEquivalentTo(getLastUsesByResearcherDto);
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsesByResearcher(researcherId),
                Times.Once);
    }

    [Fact]
    public async Task GetLastUsesByResearcherShouldThrowException()
    {
        // Assert
        var researcherId = 1;

        productRepositoryMock
            .Setup(repository => repository.GetLastUsesByResearcher(researcherId))
            .ThrowsAsync(new Exception("An error occurred while getting the last uses by researcher."));

        // Act
        Func<Task> act = () => this.controller.GetLastUsesByResearcher(researcherId);

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting the last uses by researcher.");
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsesByResearcher(researcherId),
                Times.Once);
    }

    [Fact]
    public async Task GetLastUsedProductByResearcherShouldCallRepository()
    {
        // Arrange
        var researcherId = 1;
        var getLastUsesByResearcherDto = new GetLastUsedProductDto()
        {
            ProductName = "Name",
            Quantity = 1,
            Unit = "Unit",
        };

        productRepositoryMock
            .Setup(repository => repository.GetLastUsedProductByResearcher(researcherId))
            .ReturnsAsync(getLastUsesByResearcherDto);

        // Act
        var result = await this.controller.GetLastUsedProductByResearcher(researcherId);

        // Assert
        result.Should().BeEquivalentTo(getLastUsesByResearcherDto);
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsedProductByResearcher(researcherId),
                Times.Once);
    }

    [Fact]
    public async Task GetLastUsedProductByResearcherShouldThrowException()
    {
        // Assert
        var researcherId = 1;

        productRepositoryMock
            .Setup(repository => repository.GetLastUsedProductByResearcher(researcherId))
            .ThrowsAsync(new Exception("An error occurred while getting last used product by researcher."));

        // Act
        Func<Task> act = () => this.controller.GetLastUsedProductByResearcher(researcherId);

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting last used product by researcher.");
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsedProductByResearcher(researcherId),
                Times.Once);
    }

    [Fact]
    public async Task GetLastUsesByProductShouldCallRepository()
    {
        // Arrange
        var researcherId = 1;
        var getLastUsesByResearcherDto = new List<GetLastUseByProductDto>()
        {
            new()
            {
                ResearcherName = "Name",
                Quantity = 1,
                Unit = "Unit",
            }
        };

        productRepositoryMock
            .Setup(repository => repository.GetLastUsesByProduct(researcherId))
            .ReturnsAsync(getLastUsesByResearcherDto);

        // Act
        var result = await this.controller.GetLastUsesByProduct(researcherId);

        // Assert
        result.Should().BeEquivalentTo(getLastUsesByResearcherDto);
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsesByProduct(researcherId),
                Times.Once);
    }

    [Fact]
    public async Task GetLastUsesByProductShouldThrowException()
    {
        // Assert
        var researcherId = 1;

        productRepositoryMock
            .Setup(repository => repository.GetLastUsesByProduct(researcherId))
            .ThrowsAsync(new Exception("An error occurred while getting last uses by product."));

        // Act
        Func<Task> act = () => this.controller.GetLastUsesByProduct(researcherId);

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting last uses by product.");
        productRepositoryMock
            .Verify(
                repository => repository.GetLastUsesByProduct(researcherId),
                Times.Once);
    }
}
