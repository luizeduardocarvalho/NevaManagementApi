using NevaManagement.Domain.Dtos.Invitation;

namespace NevaManagement.Domain.Interfaces.Services;

public interface ILaboratoryInvitationService
{
    Task<string> CreateInvitation(CreateInvitationDto createInvitationDto, long invitedById);
    Task<IList<GetInvitationDto>> GetLaboratoryInvitations(long laboratoryId);
    Task<InvitationDetailsDto> GetInvitationDetails(Guid token);
    Task<bool> AcceptInvitation(AcceptInvitationDto acceptInvitationDto);
    Task<bool> CancelInvitation(long invitationId, long userId);
    Task<IList<GetInvitationDto>> GetUserPendingInvitations(string email);
    Task<string> GenerateInvitationUrl(Guid token, string baseUrl);
}
