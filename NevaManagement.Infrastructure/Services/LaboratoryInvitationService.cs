using NevaManagement.Domain.Dtos.Invitation;

namespace NevaManagement.Infrastructure.Services;

public class LaboratoryInvitationService : ILaboratoryInvitationService
{
    private readonly ILaboratoryInvitationRepository invitationRepository;
    private readonly IResearcherRepository researcherRepository;
    private readonly ILaboratoryRepository laboratoryRepository;

    public LaboratoryInvitationService(
        ILaboratoryInvitationRepository invitationRepository,
        IResearcherRepository researcherRepository,
        ILaboratoryRepository laboratoryRepository)
    {
        this.invitationRepository = invitationRepository;
        this.researcherRepository = researcherRepository;
        this.laboratoryRepository = laboratoryRepository;
    }

    public async Task<string> CreateInvitation(CreateInvitationDto createInvitationDto, long invitedById)
    {
        try
        {
            // Validate laboratory exists
            if (!await laboratoryRepository.LaboratoryExists(createInvitationDto.LaboratoryId))
                throw new Exception("Laboratory not found");

            // Check if invitation already exists for this email and laboratory
            var existingInvitations = await invitationRepository.GetPendingInvitations(createInvitationDto.InviteeEmail);
            if (existingInvitations.Any(i => i.LaboratoryId == createInvitationDto.LaboratoryId))
                throw new Exception("Invitation already exists for this email and laboratory");

            var invitation = new LaboratoryInvitation
            {
                LaboratoryId = createInvitationDto.LaboratoryId,
                InviteeEmail = createInvitationDto.InviteeEmail,
                InviteeName = createInvitationDto.InviteeName,
                InvitationToken = Guid.NewGuid(),
                CreatedAt = DateTimeOffset.UtcNow,
                ExpiresAt = DateTimeOffset.UtcNow.AddDays(7), // 7 days expiration
                InvitedById = invitedById,
                Role = createInvitationDto.Role,
                IsAccepted = false,
                IsActive = true
            };

            await invitationRepository.Insert(invitation);
            await invitationRepository.SaveChanges();

            return invitation.InvitationToken.ToString();
        }
        catch
        {
            throw new Exception("An error occurred while creating the invitation.");
        }
    }

    public async Task<IList<GetInvitationDto>> GetLaboratoryInvitations(long laboratoryId)
    {
        try
        {
            return await invitationRepository.GetInvitationsByLaboratory(laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting laboratory invitations.");
        }
    }

    public async Task<InvitationDetailsDto> GetInvitationDetails(Guid token)
    {
        try
        {
            return await invitationRepository.GetInvitationByToken(token);
        }
        catch
        {
            throw new Exception("An error occurred while getting invitation details.");
        }
    }

    public async Task<bool> AcceptInvitation(AcceptInvitationDto acceptInvitationDto)
    {
        try
        {
            var invitation = await invitationRepository.GetInvitationEntityByToken(acceptInvitationDto.InvitationToken);
            
            if (invitation == null)
                return false;

            // Validate invitation
            if (invitation.IsAccepted || !invitation.IsActive || invitation.ExpiresAt < DateTimeOffset.UtcNow)
                return false;

            // Update researcher with laboratory information
            var researcher = await researcherRepository.GetById(acceptInvitationDto.UserId);
            if (researcher == null)
                return false;

            researcher.LaboratoryId = invitation.LaboratoryId;

            // Mark invitation as accepted
            invitation.IsAccepted = true;
            invitation.AcceptedAt = DateTimeOffset.UtcNow;
            invitation.AcceptedById = acceptInvitationDto.UserId;

            return await invitationRepository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while accepting the invitation.");
        }
    }

    public async Task<bool> CancelInvitation(long invitationId, long userId)
    {
        try
        {
            var invitation = await invitationRepository.GetById(invitationId);
            if (invitation == null || invitation.InvitedById != userId)
                return false;

            invitation.IsActive = false;
            return await invitationRepository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while canceling the invitation.");
        }
    }

    public async Task<IList<GetInvitationDto>> GetUserPendingInvitations(string email)
    {
        try
        {
            return await invitationRepository.GetPendingInvitations(email);
        }
        catch
        {
            throw new Exception("An error occurred while getting user pending invitations.");
        }
    }

    public Task<string> GenerateInvitationUrl(Guid token, string baseUrl)
    {
        return Task.FromResult($"{baseUrl.TrimEnd('/')}/invitation/accept/{token}");
    }
}
