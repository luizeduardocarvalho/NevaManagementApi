using NevaManagement.Domain.Models;

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
    public async Task GetEquipmentUsageCalendarShouldReturnDateList()
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

        var expectedResult = new List<EquipmentUsageCalendarDto>
        {
            new()
            {
                Month = 1,
                Days = new List<Day>
                {
                    new()
                    {
                        DayNumber = 1,
                        Appointments = new List<Appointment>
                        {
                            new()
                            {
                                StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
                                EndDate = new DateTime(2022, 1, 1, 11, 30, 0)
                            }
                        }
                    }
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
        result.Should().BeEquivalentTo(expectedResult);
    }

    [Fact]
    public async Task GetEquipmentUsageCalendarShouldThrowException()
    {
        // Arrange
        equipmentUsageRepositoryMock
            .Setup(repository => repository.GetEquipmentUsageByEquipment(1))
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => service.GetEquipmentUsageCalendar(1);

        // Ararnge
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting the equipment usage calendar.");
    }

    [Fact]
    public async Task GetEquipmentUsageHistoryShouldReturnHistory()
    {
        // Arrange
        var expectedResult = new List<GetEquipmentUsageDto>
        {
            new()
            {
                Id = 1,
            }
        };

        equipmentUsageRepositoryMock
            .Setup(repository => repository.GetEquipmentUsageByEquipment(1))
            .ReturnsAsync(expectedResult);

        // Act
        var result = await this.service.GetEquipmentUsageHistory(1);

        // Assert
        equipmentUsageRepositoryMock
            .Verify(repository => repository.GetEquipmentUsageByEquipment(1), Times.Once);
        result.Should().BeEquivalentTo(expectedResult);
    }

    [Fact]
    public async Task GetEquipmentUsageHistoryShouldReturnException()
    {
        // Ararnge
        equipmentUsageRepositoryMock
            .Setup(repository => repository.GetEquipmentUsageByEquipment(1))
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => service.GetEquipmentUsageHistory(1);

        // Ararnge
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while getting the equipment usage history.");
    }

    [Fact]
    public async Task UseEquipmentWithCorrectDatesShouldReturnTrue()
    {
        // Arrange
        var useEquipmentDto = new UseEquipmentDto()
        {
            EquipmentId = 1,
            StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
            EndDate = new DateTime(2022, 1, 1, 11, 30, 0),
            ResearcherId = 1
        };

        var equipmenHistory = new List<GetEquipmentUsageDto>
        {
            new()
            {
                Id = 1,
                StartDate = new DateTime(2022, 1, 1, 12, 30, 0),
                EndDate = new DateTime(2022, 1, 1, 13, 30, 0)
            }
        };

        equipmentUsageRepositoryMock
            .Setup(
                repository => repository.GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate))
            .ReturnsAsync(equipmenHistory);

        equipmentUsageRepositoryMock
            .Setup(repository => repository.SaveChanges())
            .ReturnsAsync(true);

        // Act
        var result = await this.service.UseEquipment(useEquipmentDto);

        // Assert
        equipmentUsageRepositoryMock
            .Verify(repository => repository.GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate), Times.Once);
        equipmentUsageRepositoryMock
            .Verify(
                repository => repository.Insert(It.IsAny<EquipmentUsage>()),
                Times.Once);
        equipmentUsageRepositoryMock
            .Verify(repository => repository.SaveChanges(), Times.Once);
        result.Should().BeTrue();
    }

    [Fact]
    public async Task UseEquipmentWithOverlappingDatesShouldReturnFalse()
    {
        // Arrange
        var useEquipmentDto = new UseEquipmentDto()
        {
            EquipmentId = 1,
            StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
            EndDate = new DateTime(2022, 1, 1, 11, 30, 0),
            ResearcherId = 1
        };

        var equipmenHistory = new List<GetEquipmentUsageDto>
        {
            new()
            {
                Id = 1,
                StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
                EndDate = new DateTime(2022, 1, 1, 13, 30, 0)
            }
        };

        equipmentUsageRepositoryMock
            .Setup(
                repository => repository.GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate))
            .ReturnsAsync(equipmenHistory);

        // Act
        var result = await this.service.UseEquipment(useEquipmentDto);

        // Assert
        equipmentUsageRepositoryMock
            .Verify(repository => repository.GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate), Times.Once);

        result.Should().BeFalse();
    }

    [Fact]
    public async Task UseEquipmentShouldThrowException()
    {
        // Arrange
        var useEquipmentDto = new UseEquipmentDto()
        {
            EquipmentId = 1,
            StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
            EndDate = new DateTime(2022, 1, 1, 11, 30, 0),
            ResearcherId = 1
        };

        var equipmenHistory = new List<GetEquipmentUsageDto>
        {
            new()
            {
                Id = 1,
                StartDate = new DateTime(2022, 1, 1, 12, 30, 0),
                EndDate = new DateTime(2022, 1, 1, 13, 30, 0)
            }
        };

        equipmentUsageRepositoryMock
            .Setup(
                repository => repository.GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate))
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => this.service.UseEquipment(useEquipmentDto);

        // Assert
        await act
            .Should()
            .ThrowAsync<Exception>()
            .WithMessage("An error occurred while trying to use the equipment");
    }
}
