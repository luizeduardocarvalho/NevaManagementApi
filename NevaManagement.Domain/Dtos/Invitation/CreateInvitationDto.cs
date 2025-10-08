namespace NevaManagement.Domain.Dtos.Invitation;

public class CreateInvitationDto
{
    public long LaboratoryId { get; set; }
    public string InviteeEmail { get; set; }
    public string InviteeName { get; set; }
    public string Role { get; set; } // "Admin" or "Member"
}