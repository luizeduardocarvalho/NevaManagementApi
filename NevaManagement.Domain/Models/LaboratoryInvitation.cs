namespace NevaManagement.Domain.Models;

[Table("LaboratoryInvitation")]
public class LaboratoryInvitation : BaseEntity
{
    [ForeignKey("Laboratory_Id")]
    public Laboratory Laboratory { get; set; }
    public long LaboratoryId { get; set; }

    [MaxLength(100)]
    public string InviteeEmail { get; set; }

    [MaxLength(200)]
    public string InviteeName { get; set; }

    public Guid InvitationToken { get; set; }

    public DateTimeOffset CreatedAt { get; set; }

    public DateTimeOffset ExpiresAt { get; set; }

    public bool IsAccepted { get; set; }

    public DateTimeOffset? AcceptedAt { get; set; }

    [ForeignKey("InvitedBy_Id")]
    public Researcher InvitedBy { get; set; }
    public long InvitedById { get; set; }

    [ForeignKey("AcceptedBy_Id")]
    public Researcher AcceptedBy { get; set; }
    public long? AcceptedById { get; set; }

    [MaxLength(50)]
    public string Role { get; set; } // "Admin", "Member"

    public bool IsActive { get; set; }
}