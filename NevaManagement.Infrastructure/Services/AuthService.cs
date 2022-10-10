namespace NevaManagement.Infrastructure.Services;

public class AuthService : IAuthService
{
    private readonly IEncryptService encryptService;
    private readonly IResearcherService researcherService;
    private readonly ILogger<AuthService> logger;

    public AuthService(
        IEncryptService encryptService,
        IResearcherService researcherService,
        ILogger<AuthService> logger)
    {
        this.encryptService = encryptService;
        this.researcherService = researcherService;
        this.logger = logger;
    }

    async Task IAuthService.ChangePassword(ChangePasswordDto changePasswordDto)
    {
        try
        {
            var encryptedPassword = await this.encryptService
                .Encrypt(changePasswordDto.OldPassword);
            var user = await this.researcherService.GetByEmailAndPassword(
                changePasswordDto.Email,
                encryptedPassword);
            
            if (user != null)
            {
                var entity = await this.researcherService.GetById(user.Id);

                var newEncryptedPassword = await this.encryptService
                    .Encrypt(changePasswordDto.NewPassword);

                entity.Password = newEncryptedPassword;

                await this.researcherService.Save().ConfigureAwait(false);
            }
            else
            {
                throw new WrongPasswordException("Wrong password");
            }
        }
        catch(WrongPasswordException e)
        {
            throw new WrongPasswordException(e.Message);
        }
        catch
        {
            throw new Exception("An error occurred while changing the password.");
        }
    }
}