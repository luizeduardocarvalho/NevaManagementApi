namespace NevaManagement.Tests.Controllers;

public sealed class ContainerControllerTests
{
    private readonly Mock<IContainerService> containerService;
    private readonly ContainerController containerController;

    public ContainerControllerTests()
    {
        this.containerService = new Mock<IContainerService>();
        this.containerController = new ContainerController(this.containerService.Object);
    }

    [Fact]
    public async Task GetContainers_ShouldReturn200Ok()
    {
        // Arrange
        var expectedContainerList = new List<GetSimpleContainerDto>()
        {
            new()
            {
                Id = 1,
                Name = "ContainerName"
            }
        };

        containerService.Setup(service => service.GetContainers()).ReturnsAsync(expectedContainerList);

        // Act
        var result = await containerController.GetContainers();

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        containerService.Verify(service => service.GetContainers(), Times.Once);
    }

    [Fact]
    public async Task AddContainers_WithValidInput_ShouldReturn201Created()
    {
        // Arrange
        var newContainer = new AddContainerDto()
        {
            Name = "ContainerName",
            ResearcherId = 1,
            OrganismId = 1
        };

        var expectedResult = new ObjectResult($"Successfully created {newContainer.Name}.")
        {
            StatusCode = 201
        };

        containerService.Setup(service => service.AddContainer(newContainer)).ReturnsAsync(true);

        // Act
        var result = await containerController.AddContainer(newContainer);

        // Assert
        result.Should().BeEquivalentTo(expectedResult);
        containerService.Verify(service => service.AddContainer(newContainer), Times.Once);

    }

    [Fact]
    public async Task AddContainers_WithErrorInModelState_ShouldReturn400BadRequest()
    {
        // Arrange
        var newContainer = new AddContainerDto()
        {
            Name = "ContainerName",
            OrganismId = 1
        };

        containerController.ModelState.AddModelError("Error", "Error");

        // Act
        var result = await containerController.AddContainer(newContainer);

        // Assert
        result.Should().BeOfType<BadRequestResult>();
    }

    [Theory]
    [InlineData(1)]
    [InlineData(100)]
    [InlineData(10000)]
    public async Task GetChildrenContainers_WithValidId_ShouldReturn200Ok(long id)
    {
        // Arrange
        var expectedContainerList = new List<GetSimpleContainerDto>()
        {
            new()
            {
                Id = 1,
                Name = "ContainerName"
            }
        };

        containerService.Setup(service => service.GetChildrenContainers(id)).ReturnsAsync(expectedContainerList);

        // Act
        var result = await containerController.GetChildrenContainers(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        containerService.Verify(service => service.GetChildrenContainers(id), Times.Once);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-10)]
    public async Task GetChildrenContainers_WithInValidId_ShouldReturn400BadRequest(long id)
    {
        // Arrange
        containerController.ModelState.AddModelError("Error", "Error");

        // Act
        var result = await containerController.GetChildrenContainers(id);

        // Assert
        result.Should().BeOfType<BadRequestResult>();
    }

    [Theory]
    [InlineData(1)]
    [InlineData(100)]
    [InlineData(10000)]
    public async Task GetDetailedContainer_WithValidId_ShouldReturn200Ok(long id)
    {
        // Arrange
        var expectedContainer = new GetDetailedContainerDto()
        {
            Name = "ContainerName"
        };

        containerService.Setup(service => service.GetDetailedContainer(id)).ReturnsAsync(expectedContainer);

        // Act
        var result = await containerController.GetDetailedContainer(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        containerService.Verify(service => service.GetDetailedContainer(id), Times.Once);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-10)]
    public async Task GetDetailedContainer_WithInValidId_ShouldReturn400BadRequest(long id)
    {
        // Arrange
        containerController.ModelState.AddModelError("Error", "Error");

        // Act
        var result = await containerController.GetDetailedContainer(id);

        // Assert
        result.Should().BeOfType<BadRequestResult>();
    }
}
