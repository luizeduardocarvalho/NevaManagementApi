using NevaManagement.Domain.Dtos.Invitation;
using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Domain.Dtos.Onboarding;

public class OnboardingStatusDto
{
    public bool IsOnboardingComplete { get; set; }
    public bool HasLaboratory { get; set; }
    public GetSimpleLaboratoryDto Laboratory { get; set; }
    public IList<GetInvitationDto> PendingInvitations { get; set; }
}
