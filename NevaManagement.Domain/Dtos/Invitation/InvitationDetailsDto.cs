namespace NevaManagement.Domain.Dtos.Invitation;

public class InvitationDetailsDto
{
    public long LaboratoryId { get; set; }
    public string LaboratoryName { get; set; }
    public string LaboratoryDescription { get; set; }
    public string InviteeName { get; set; }
    public string Role { get; set; }
    public string InvitedByName { get; set; }
    public DateTimeOffset ExpiresAt { get; set; }
    public bool IsExpired { get; set; }
    public bool IsAccepted { get; set; }
}