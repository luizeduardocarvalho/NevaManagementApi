using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Domain.Dtos.Onboarding;

public class OnboardingChoiceDto
{
    public string Choice { get; set; } // "CreateLab" or "JoinLab"
    public CreateLaboratoryDto CreateLaboratory { get; set; } // For creating new lab
    public Guid? InvitationToken { get; set; } // For joining existing lab
    public long UserId { get; set; }
}
