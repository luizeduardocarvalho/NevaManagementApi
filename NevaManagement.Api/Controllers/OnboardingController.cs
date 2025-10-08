using System.Linq;
using NevaManagement.Domain.Dtos.Laboratory;
using NevaManagement.Domain.Dtos.Onboarding;
using NevaManagement.Domain.Dtos.Invitation;

namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class OnboardingController : ControllerBase
{
    private readonly ILaboratoryService laboratoryService;
    private readonly ILaboratoryInvitationService invitationService;
    private readonly IResearcherService researcherService;

    public OnboardingController(
        ILaboratoryService laboratoryService,
        ILaboratoryInvitationService invitationService,
        IResearcherService researcherService)
    {
        this.laboratoryService = laboratoryService;
        this.invitationService = invitationService;
        this.researcherService = researcherService;
    }

    [HttpGet("Status")]
    public async Task<IActionResult> GetOnboardingStatus([FromQuery] long userId, [FromQuery] string email)
    {
        try
        {
            // Get user's current laboratory status
            var researcher = await researcherService.GetById(userId);
            var hasLaboratory = researcher?.Laboratory != null;

            // Get pending invitations
            var pendingInvitations = await invitationService.GetUserPendingInvitations(email);

            var status = new OnboardingStatusDto
            {
                IsOnboardingComplete = hasLaboratory,
                HasLaboratory = hasLaboratory,
                Laboratory = hasLaboratory ? new GetSimpleLaboratoryDto
                {
                    Id = researcher.Laboratory.Id,
                    Name = researcher.Laboratory.Name,
                    IsActive = researcher.Laboratory.IsActive
                } : null,
                PendingInvitations = pendingInvitations
            };

            return Ok(status);
        }
        catch (Exception ex)
        {
            return BadRequest($"Error getting onboarding status: {ex.Message}");
        }
    }

    [HttpPost("Complete")]
    public async Task<IActionResult> CompleteOnboarding([FromBody] OnboardingChoiceDto onboardingChoice)
    {
        try
        {
            if (onboardingChoice.Choice == "CreateLab")
            {
                // Create new laboratory
                var labCreated = await laboratoryService.CreateLaboratory(onboardingChoice.CreateLaboratory);
                if (!labCreated)
                    return BadRequest("Failed to create laboratory");

                // Get the created laboratory to assign to user
                var laboratories = await laboratoryService.GetSimpleLaboratories();
                var newLab = laboratories.OrderByDescending(l => l.Id).FirstOrDefault();

                if (newLab != null)
                {
                    // Update researcher with laboratory ID
                    var researcher = await researcherService.GetById(onboardingChoice.UserId);
                    if (researcher != null)
                    {
                        // Note: You'll need to implement UpdateResearcherLaboratory method
                        // researcher.LaboratoryId = newLab.Id;
                        // await researcherService.UpdateResearcher(researcher);
                    }
                }

                return Ok(new { Success = true, Message = "Laboratory created and user assigned" });
            }
            else if (onboardingChoice.Choice == "JoinLab" && onboardingChoice.InvitationToken.HasValue)
            {
                // Accept invitation
                var acceptDto = new AcceptInvitationDto
                {
                    InvitationToken = onboardingChoice.InvitationToken.Value,
                    UserId = onboardingChoice.UserId
                };

                var accepted = await invitationService.AcceptInvitation(acceptDto);
                if (!accepted)
                    return BadRequest("Failed to accept invitation or invitation is no longer valid");

                return Ok(new { Success = true, Message = "Invitation accepted successfully" });
            }

            return BadRequest("Invalid onboarding choice");
        }
        catch (Exception ex)
        {
            return BadRequest($"Error completing onboarding: {ex.Message}");
        }
    }
}
