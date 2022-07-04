using NevaManagement.Domain.Dtos.Researcher;

namespace NevaManagement.Tests.Controllers;

public sealed class AuthControllerTests
{
    private readonly Mock<IEncryptService> encryptServiceMock;
    private readonly Mock<IResearcherService> researcherServiceMock;
    private readonly Mock<ITokenService> tokenServiceMock;
    private readonly Mock<IAuthService> authServiceMock;
    private readonly AuthController controller;

    public AuthControllerTests()
    {
        encryptServiceMock = new Mock<IEncryptService>();
        researcherServiceMock = new Mock<IResearcherService>();
        tokenServiceMock = new Mock<ITokenService>();
        authServiceMock = new Mock<IAuthService>();

        controller = new AuthController(
            encryptServiceMock.Object,
            researcherServiceMock.Object,
            tokenServiceMock.Object,
            authServiceMock.Object);
    }

    [Fact]
    public async Task AuthenticateAsyncShouldReturnOkWithObject()
    {
        // Arrange
        var login = new LoginDto()
        {
            Email = "email@email.com",
            Password = "Password"
        };

        const string encryptedPassword = "EncryptedPassword";
        const string token = "Token";

        var researcher = new GetDetailedResearcherDto()
        {
            Id = 1,
            Email = "email@email.com",
            Password = encryptedPassword
        };

        var expectedResult = new LoggedUserDto()
        {
            Researcher = researcher,
            Token = token
        };

        encryptServiceMock.Setup(s => s.Encrypt(login.Password)).ReturnsAsync(encryptedPassword);
        researcherServiceMock
            .Setup(s => s.GetByEmailAndPassword(login.Email, encryptedPassword))
            .ReturnsAsync(researcher);
        tokenServiceMock
            .Setup(s => s.GenerateToken(researcher))
            .ReturnsAsync(token);

        // Act
        var result = await controller.AuthenticateAsync(login);

        // Assert
        result.Should().BeOfType<OkObjectResult>()
            .Which.Value.Should().BeEquivalentTo(expectedResult);
        encryptServiceMock.Verify(s => s.Encrypt(login.Password), Times.Once);
        tokenServiceMock.Verify(s => s.GenerateToken(researcher), Times.Once);
        researcherServiceMock
            .Verify(
                s => s.GetByEmailAndPassword(
                    login.Email,
                    encryptedPassword),
                Times.Once);
    }

    [Fact]
    public async Task AuthenticateAsyncWithNullUserShouldReturnNotFound()
    {
        // Arrange
        var login = new LoginDto()
        {
            Email = "email@email.com",
            Password = "Password"
        };

        const string encryptedPassword = "EncryptedPassword";

        encryptServiceMock.Setup(s => s.Encrypt(login.Password)).ReturnsAsync(encryptedPassword);
        researcherServiceMock
            .Setup(s => s.GetByEmailAndPassword(login.Email, encryptedPassword))
            .ReturnsAsync((GetDetailedResearcherDto)null);

        // Act
        var result = await this.controller.AuthenticateAsync(login);

        // Assert
        result.Should().BeOfType<NotFoundObjectResult>();
    }

    [Fact]
    public async Task RegisterShouldReturnOk()
    {
        // Arrange
        var register = new RegisterDto()
        {
            Email = "email@email.com",
            Password = "Password"
        };

        const string encryptedPassword = "EncryptedPassword";

        encryptServiceMock.Setup(s => s.Encrypt(register.Password)).ReturnsAsync(encryptedPassword);
        authServiceMock.Setup(s => s.Register(register)).ReturnsAsync(true);

        // Act
        var result = await this.controller.Register(register);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        encryptServiceMock.Verify(s => s.Encrypt(register.Password), Times.Once);
        authServiceMock.Verify(s => s.Register(register), Times.Once);
    }

    [Fact]
    public async Task RegisterWithoutPasswordShouldReturnCreateWithDefaultPassword()
    {
        // Arrange
        var register = new RegisterDto()
        {
            Email = "email@email.com"
        };

        const string defaultPassword = "123456";
        const string encryptedPassword = "EncryptedPassword";

        encryptServiceMock.Setup(s => s.Encrypt(defaultPassword)).ReturnsAsync(encryptedPassword);
        authServiceMock.Setup(s => s.Register(register)).ReturnsAsync(true);

        // Act
        var result = await this.controller.Register(register);

        // Assert
        result.Should().BeOfType<OkObjectResult>();
        encryptServiceMock.Verify(s => s.Encrypt(defaultPassword), Times.Once);
        authServiceMock.Verify(s => s.Register(register), Times.Once);
    }

    [Fact]
    public async Task RegisterWithErrorOnSaveShouldReturnInternalServerError()
    {
        // Arrange
        var register = new RegisterDto()
        {
            Email = "email@email.com",
            Password = "Password"
        };

        const string encryptedPassword = "EncryptedPassword";

        encryptServiceMock.Setup(s => s.Encrypt(register.Password)).ReturnsAsync(encryptedPassword);
        authServiceMock.Setup(s => s.Register(register)).ReturnsAsync(false);

        // Act
        var result = await this.controller.Register(register);

        // Assert
        result.Should().BeOfType<StatusCodeResult>().Which.StatusCode.Should().Be(500);
        encryptServiceMock.Verify(s => s.Encrypt(register.Password), Times.Once);
        authServiceMock.Verify(s => s.Register(register), Times.Once);
    }
}
