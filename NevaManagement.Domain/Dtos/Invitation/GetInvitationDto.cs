namespace NevaManagement.Domain.Dtos.Invitation;

public class GetInvitationDto
{
    public long Id { get; set; }
    public long LaboratoryId { get; set; }
    public string LaboratoryName { get; set; }
    public string InviteeEmail { get; set; }
    public string InviteeName { get; set; }
    public string Role { get; set; }
    public DateTimeOffset CreatedAt { get; set; }
    public DateTimeOffset ExpiresAt { get; set; }
    public bool IsAccepted { get; set; }
    public DateTimeOffset? AcceptedAt { get; set; }
    public string InvitedByName { get; set; }
    public bool IsExpired { get; set; }
}