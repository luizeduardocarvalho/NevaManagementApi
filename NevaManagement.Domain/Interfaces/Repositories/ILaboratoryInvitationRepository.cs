using NevaManagement.Domain.Dtos.Invitation;

namespace NevaManagement.Domain.Interfaces.Repositories;

public interface ILaboratoryInvitationRepository : IBaseRepository<LaboratoryInvitation>
{
    Task<IList<GetInvitationDto>> GetInvitationsByLaboratory(long laboratoryId);
    Task<InvitationDetailsDto> GetInvitationByToken(Guid token);
    Task<LaboratoryInvitation> GetInvitationEntityByToken(Guid token);
    Task<bool> IsInvitationValid(Guid token);
    Task<IList<GetInvitationDto>> GetPendingInvitations(string email);
}
