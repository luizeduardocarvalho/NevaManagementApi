using NevaManagement.Domain.Dtos.Organism;

namespace NevaManagement.Tests.Controllers;

public sealed class OrganismControllerTests
{
    private readonly Mock<IOrganismService> organismServiceMock;
    private readonly OrganismController controller;

    public OrganismControllerTests()
    {
        organismServiceMock = new Mock<IOrganismService>();

        controller = new OrganismController(organismServiceMock.Object);
    }

    [Fact]
    public async Task GetOrganismsShouldReturnOk()
    {
        // Arrange
        var organisms = new List<GetOrganismDto>()
        {
            new()
            {
                Id = 1,
                Name = "Organism"
            }
        };

        organismServiceMock.Setup(s => s.GetOrganisms()).ReturnsAsync(organisms);

        // Act
        var result = await this.controller.GetOrganisms();

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        organismServiceMock.Verify(s => s.GetOrganisms(), Times.Once);
    }

    [Fact]
    public async Task AddOrganismShouldReturnOk()
    {
        // Arrange
        var organism = new AddOrganismDto()
        {
            Name = "Organism"
        };

        organismServiceMock.Setup(s => s.AddOrganism(organism)).ReturnsAsync(true);

        // Act
        var result = await this.controller.AddOrganism(organism);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        organismServiceMock.Verify(s => s.AddOrganism(organism), Times.Once);
    }

    [Fact]
    public async Task AddOrganismWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var organism = new AddOrganismDto()
        {
            Name = "Organism"
        };

        organismServiceMock.Setup(s => s.AddOrganism(organism)).ReturnsAsync(false);

        // Act
        var result = await this.controller.AddOrganism(organism);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        organismServiceMock.Verify(s => s.AddOrganism(organism), Times.Once);
    }

    [Fact]
    public async Task EditOrganismShouldReturnOk()
    {
        // Arrange
        var organism = new EditOrganismDto()
        {
            Id = 1,
            Name = "Organism"
        };

        organismServiceMock.Setup(s => s.EditOrganism(organism)).ReturnsAsync(true);

        // Act
        var result = await this.controller.EditOrganism(organism);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        organismServiceMock.Verify(s => s.EditOrganism(organism), Times.Once);
    }

    [Fact]
    public async Task EditOrganismWithErrorOnSaveShouldReturnOk()
    {
        // Arrange
        var organism = new EditOrganismDto()
        {
            Id = 1,
            Name = "Organism"
        };

        organismServiceMock.Setup(s => s.EditOrganism(organism)).ReturnsAsync(false);

        // Act
        var result = await this.controller.EditOrganism(organism);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        organismServiceMock.Verify(s => s.EditOrganism(organism), Times.Once);
    }

    [Fact]
    public async Task GetOrganismByIdShouldReturnOk()
    {
        // Arrange
        var id = 1;
        var organism = new GetDetailedOrganismDto()
        {
            Name = "Organism"
        };

        organismServiceMock.Setup(s => s.GetOrganismById(id)).ReturnsAsync(organism);

        // Act
        var result = await this.controller.GetOrganismById(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        organismServiceMock.Verify(s => s.GetOrganismById(id), Times.Once);
    }
}
