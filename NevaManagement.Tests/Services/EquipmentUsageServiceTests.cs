namespace NevaManagement.Tests.Services;

public sealed class EquipmentUsageServiceTests
{
    private readonly Mock<IEquipmentUsageRepository> equipmentUsageRepositoryMock;
    private readonly EquipmentUsageService service;

    public EquipmentUsageServiceTests()
    {
        equipmentUsageRepositoryMock = new Mock<IEquipmentUsageRepository>();
        service = new EquipmentUsageService(equipmentUsageRepositoryMock.Object);
    }

    [Fact]
    public async Task GetEquipmentUsageCalendarShouldReturnObject()
    {
        // Arrange
        var equipmentUsage = new List<GetEquipmentUsageDto>()
        {
            new()
            {
                Id = 1,
                Description = "Description",
                StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
                EndDate = new DateTime(2022, 1, 1, 11, 30, 0)
            }
        };

        var expectedResult = new
        {
            Month = 1,
            Days = new
            {
                Day = 1,
                Times = new
                {
                    StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
                    EndDate = new DateTime(2022, 1, 1, 11, 30, 0)
                }
            }
        };

        equipmentUsageRepositoryMock
            .Setup(repository => repository.GetEquipmentUsageByEquipment(1))
            .ReturnsAsync(equipmentUsage);

        // Act
        var result = await service.GetEquipmentUsageCalendar(1);

        // Assert
        equipmentUsageRepositoryMock
            .Verify(repository => repository.GetEquipmentUsageByEquipment(1), Times.Once);

    }
}
