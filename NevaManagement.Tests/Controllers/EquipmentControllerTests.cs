
namespace NevaManagement.Tests.Controllers;
public sealed class EquipmentControllerTests
{
    private readonly Mock<IEquipmentService> equipmentServiceMock;
    private readonly EquipmentController controller;

    public EquipmentControllerTests()
    {
        equipmentServiceMock = new Mock<IEquipmentService>();

        controller = new EquipmentController(equipmentServiceMock.Object);
    }

    [Fact]
    public async Task GetEquipmentsShouldReturnOk()
    {
        // Arrange
        var equipments = new List<GetSimpleEquipmentDto>
        {
            new()
            {
                Id = 1,
                Name = "Equipment"
            }
        };

        equipmentServiceMock
            .Setup(s => s.GetEquipments())
            .ReturnsAsync(equipments);

        // Act
        var result = await this.controller.GetEquipments();

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        equipmentServiceMock.Verify(s => s.GetEquipments(), Times.Once);
    }

    [Fact]
    public async Task AddEquipmentWithCorrectParameterShouldReturnOk()
    {
        // Arrange
        var equipment = new AddEquipmentDto()
        {
            Name = "Equipment",
            Description = "Description"
        };

        equipmentServiceMock
            .Setup(s => s.AddEquipment(equipment))
            .ReturnsAsync(true);

        // Act
        var result = await controller.AddEquipment(equipment);

        // Arrange
        result.Should().BeOfType<OkObjectResult>();
        equipmentServiceMock.Verify(s => s.AddEquipment(equipment), Times.Once);
    }

    [Fact]
    public async Task AddEquipmentWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var equipment = new AddEquipmentDto()
        {
            Name = "Equipment",
            Description = "Description"
        };

        equipmentServiceMock
            .Setup(s => s.AddEquipment(equipment))
            .ReturnsAsync(false);

        // Act
        var result = await controller.AddEquipment(equipment);

        // Arrange
        result.Should().BeOfType<ObjectResult>().Which.StatusCode.Should().Be(500);
        equipmentServiceMock.Verify(s => s.AddEquipment(equipment), Times.Once);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-10)]
    public async Task GetDetailedEquipmentWithInvalidIdShouldReturnBadRequest(long id)
    {
        // Arrange
        controller.ModelState.AddModelError("Error", "Error");

        // Act
        var result = await controller.GetDetailedEquipment(id);

        // Assert
        result.Should().BeOfType<BadRequestResult>();
    }

    [Fact]
    public async Task GetDetailedEquipmentWithCorrectIdShouldReturnOk()
    {
        // Arrange
        var id = 1;
        var equipment = new GetDetailedEquipmentDto()
        {
            Id = id,
            Name = "Equipment"
        };

        equipmentServiceMock
            .Setup(s => s.GetDetailedEquipment(id))
            .ReturnsAsync(equipment);

        // Act
        var result = await controller.GetDetailedEquipment(id);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        equipmentServiceMock.Verify(s => s.GetDetailedEquipment(id), Times.Once);
    }

    [Fact]
    public async Task EditEquipmentWithValidParameterShouldReturnOk()
    {
        // Arrange
        var equipment = new EditEquipmentDto()
        {
            Id = 1,
            Name = "Equipment"
        };

        equipmentServiceMock
            .Setup(s => s.EditEquipment(equipment))
            .ReturnsAsync(true);

        // Act
        var result = await controller.EditEquipment(equipment);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        equipmentServiceMock.Verify(s => s.EditEquipment(equipment), Times.Once);
    }

    [Fact]
    public async Task EditEquipmentWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var equipment = new EditEquipmentDto()
        {
            Id = 1,
            Name = "Equipment"
        };

        equipmentServiceMock
            .Setup(s => s.EditEquipment(equipment))
            .ReturnsAsync(false);

        // Act
        var result = await controller.EditEquipment(equipment);

        // Assert
        result.Should().BeOfType<ObjectResult>().Which.StatusCode.Should().Be(500);
        equipmentServiceMock.Verify(s => s.EditEquipment(equipment), Times.Once);
    }
}
