namespace NevaManagement.Domain.Dtos.Invitation;

public class AcceptInvitationDto
{
    public Guid InvitationToken { get; set; }
    public long UserId { get; set; } // The user accepting the invitation
}