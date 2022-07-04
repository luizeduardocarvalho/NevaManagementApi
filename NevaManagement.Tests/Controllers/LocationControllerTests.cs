using NevaManagement.Domain.Dtos.Location;

namespace NevaManagement.Tests.Controllers;

public sealed class LocationControllerTests
{
    private readonly Mock<ILocationService> locationServiceMock;
    private readonly LocationController controller;

    public LocationControllerTests()
    {
        locationServiceMock = new Mock<ILocationService>();

        controller = new LocationController(locationServiceMock.Object);
    }

    [Fact]
    public async Task GetLocationsShouldReturnOk()
    {
        // Arrange
        var locations = new List<GetLocationDto>()
        {
            new()
            {
                Id = 1,
                Name = "Location"
            }
        };

        locationServiceMock.Setup(s => s.GetLocations()).ReturnsAsync(locations);

        // Act
        var result = await this.controller.GetLocations();

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        locationServiceMock.Verify(s => s.GetLocations(), Times.Once);
    }

    [Fact]
    public async Task AddLocationShouldReturnOk()
    {
        // Arrange
        var location = new AddLocationDto()
        {
            Name = "Location"
        };

        locationServiceMock.Setup(s => s.AddLocation(location)).ReturnsAsync(true);

        // Act
        var result = await this.controller.AddLocation(location);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        locationServiceMock.Verify(s => s.AddLocation(location), Times.Once);
    }

    [Fact]
    public async Task AddLocationWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var location = new AddLocationDto()
        {
            Name = "Location"
        };

        locationServiceMock.Setup(s => s.AddLocation(location)).ReturnsAsync(false);

        // Act
        var result = await this.controller.AddLocation(location);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        locationServiceMock.Verify(s => s.AddLocation(location), Times.Once);
    }

    [Fact]
    public async Task EditLocationShouldReturnOk()
    {
        // Arrange
        var location = new EditLocationDto()
        {
            Id = 1,
            Name = "Location"
        };

        locationServiceMock.Setup(s => s.EditLocation(location)).ReturnsAsync(true);

        // Act
        var result = await this.controller.EditLocation(location);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        locationServiceMock.Verify(s => s.EditLocation(location), Times.Once);
    }

    [Fact]
    public async Task EditLocationWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var location = new EditLocationDto()
        {
            Id = 1,
            Name = "Location"
        };

        locationServiceMock.Setup(s => s.EditLocation(location)).ReturnsAsync(false);

        // Act
        var result = await this.controller.EditLocation(location);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        locationServiceMock.Verify(s => s.EditLocation(location), Times.Once);
    }

    [Fact]
    public async Task GetLocationByIdShouldReturnOk()
    {
        // Arrange
        var id = 1;
        var location = new GetDetailedLocationDto()
        {
            Name = "Location"
        };

        locationServiceMock.Setup(s => s.GetLocationById(id)).ReturnsAsync(location);

        // Act
        var result = await this.controller.GetLocationById(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        locationServiceMock.Verify(s => s.GetLocationById(id), Times.Once);
    }
}
