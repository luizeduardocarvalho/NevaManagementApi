namespace NevaManagement.Tests.Services;

public sealed class ContainerServiceTests
{
    private readonly Mock<IContainerRepository> containerRepositoryMock;
    private readonly Mock<IResearcherRepository> researcherRepositoryMock;
    private readonly Mock<IOrganismRepository> organismRepositoryMock;
    private readonly ContainerService containerService;

    public ContainerServiceTests()
    {
        containerRepositoryMock = new Mock<IContainerRepository>();
        researcherRepositoryMock = new Mock<IResearcherRepository>();
        organismRepositoryMock = new Mock<IOrganismRepository>();

        containerService = new ContainerService(
            containerRepositoryMock.Object,
            researcherRepositoryMock.Object,
            organismRepositoryMock.Object);
    }

    [Fact]
    public async Task GetContainers_ShouldThrowException()
    {
        // Arrange
        var simpleContainerDtoListExpected = new List<GetSimpleContainerDto>
            {
                new()
                {
                    Id = 1,
                    Name = "Name"
                }
            };

        containerRepositoryMock
            .Setup(repository => repository.GetContainers())
            .ReturnsAsync(simpleContainerDtoListExpected);

        // Act
        await containerService.GetContainers();

        // Arrange
        containerRepositoryMock.Verify(repository => repository.GetContainers(), Times.Once);
    }

    [Fact]
    public async Task GetContainers_WithRepositoryException_ShouldThrowException()
    {
        // Arrange
        containerRepositoryMock
            .Setup(repository => repository.GetContainers())
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => containerService.GetContainers();

        // Arrange
        await act.Should().ThrowAsync<Exception>().WithMessage("An error occurred while getting all containers.");
        containerRepositoryMock.Verify(repository => repository.GetContainers(), Times.Once);
    }

    [Theory]
    [InlineData(1)]
    [InlineData(100)]
    public async Task GetChildrenContainers_ShouldThrowException(long id)
    {
        // Arrange
        var simpleContainerDtoListExpected = new List<GetSimpleContainerDto>
            {
                new()
                {
                    Id = 1,
                    Name = "Name"
                }
            };

        containerRepositoryMock
            .Setup(repository => repository.GetChildrenContainers(id))
            .ReturnsAsync(simpleContainerDtoListExpected);

        // Act
        await containerService.GetChildrenContainers(id);

        // Arrange
        containerRepositoryMock.Verify(repository => repository.GetChildrenContainers(id), Times.Once);
    }

    [Theory]
    [InlineData(0)]
    [InlineData(-10)]
    public async Task GetChildrenContainers_WithRepositoryException_ShouldThrowException(long id)
    {
        // Arrange
        containerRepositoryMock
            .Setup(repository => repository.GetChildrenContainers(id))
            .ThrowsAsync(new Exception());

        // Act
        Func<Task> act = () => containerService.GetChildrenContainers(id);

        // Arrange
        await act.Should().ThrowAsync<Exception>().WithMessage("An error occurred while getting the children for the container.");
        containerRepositoryMock.Verify(repository => repository.GetChildrenContainers(id), Times.Once);
    }
}
