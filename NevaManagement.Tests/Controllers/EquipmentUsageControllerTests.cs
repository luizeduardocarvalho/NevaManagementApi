namespace NevaManagement.Tests.Controllers;

public class EquipmentUsageControllerTests
{
    private readonly Mock<IEquipmentUsageService> equipmentUsageServiceMock;
    private readonly EquipmentUsageController controller;

    public EquipmentUsageControllerTests()
    {
        equipmentUsageServiceMock = new Mock<IEquipmentUsageService>();
        controller = new EquipmentUsageController(equipmentUsageServiceMock.Object);
    }


    [Fact]
    public async Task GetEquipmentUsageCalendarShouldReturnOk()
    {
        // Arrange
        equipmentUsageServiceMock
            .Setup(service => service.GetEquipmentUsageCalendar(1))
            .ReturnsAsync(It.IsAny<IList<EquipmentUsageCalendarDto>>);

        // Act
        var result = await controller.GetEquipmentUsageCalendar(1);

        // Arrange
        result.Should().BeOfType<OkObjectResult>();
        equipmentUsageServiceMock
            .Verify(service => service.GetEquipmentUsageCalendar(1), Times.Once);
    }

    [Fact]
    public async Task GetEquipmentUsageHistoryShouldReturnOk()
    {
        // Arrange
        var expectedResult = new List<GetEquipmentUsageDto>()
        {
            new()
            {
                Id = 1,
                Description = "Description"
            }
        };

        equipmentUsageServiceMock
            .Setup(service => service.GetEquipmentUsageHistory(1))
            .ReturnsAsync(expectedResult);

        // Act
        var result = await controller.GetEquipmentUsageHistory(1);

        // Arrange
        result.Should().BeOfType<OkObjectResult>();
        equipmentUsageServiceMock
            .Verify(service => service.GetEquipmentUsageHistory(1), Times.Once);
    }

    [Fact]
    public async Task UseEquipmentShouldReturnOk()
    {
        // Arrange
        var useEquipment = new UseEquipmentDto()
        {
            ResearcherId = 1,
            EquipmentId = 1
        };

        equipmentUsageServiceMock
            .Setup(service => service.UseEquipment(useEquipment))
            .ReturnsAsync(It.IsAny<bool>());

        // Act
        var result = await controller.UseEquipment(useEquipment);

        // Arrange
        result.Should().BeOfType<OkObjectResult>();
        equipmentUsageServiceMock
            .Verify(service => service.UseEquipment(useEquipment), Times.Once);
    }
}
