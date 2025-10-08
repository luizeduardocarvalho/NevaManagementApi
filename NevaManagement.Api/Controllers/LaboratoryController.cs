using NevaManagement.Domain.Dtos.Invitation;
using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Api.Controllers;

[ApiController]
[Route("[controller]")]
[Produces("application/json")]
public class LaboratoryController : ControllerBase
{
    private readonly ILaboratoryService service;
    private readonly ILaboratoryInvitationService invitationService;

    public LaboratoryController(ILaboratoryService service, ILaboratoryInvitationService invitationService)
    {
        this.service = service;
        this.invitationService = invitationService;
    }

    [HttpGet("GetAll")]
    public async Task<IActionResult> GetAll()
    {
        var laboratories = await service.GetAllLaboratories();
        return Ok(laboratories);
    }

    [HttpGet("GetSimple")]
    public async Task<IActionResult> GetSimple()
    {
        var laboratories = await service.GetSimpleLaboratories();
        return Ok(laboratories);
    }

    [HttpGet("GetById")]
    public async Task<IActionResult> GetById([FromQuery] long id)
    {
        var laboratory = await service.GetLaboratoryById(id);
        
        if (laboratory == null)
            return NotFound();

        return Ok(laboratory);
    }

    [HttpPost("Create")]
    public async Task<IActionResult> Create([FromBody] CreateLaboratoryDto createLaboratoryDto)
    {
        var result = await service.CreateLaboratory(createLaboratoryDto);
        
        if (!result)
            return BadRequest("Failed to create laboratory");

        return Ok(result);
    }

    [HttpPut("Update")]
    public async Task<IActionResult> Update([FromBody] EditLaboratoryDto editLaboratoryDto)
    {
        var result = await service.UpdateLaboratory(editLaboratoryDto);
        
        if (!result)
            return NotFound("Laboratory not found");

        return Ok(result);
    }

    [HttpDelete("Delete")]
    public async Task<IActionResult> Delete([FromQuery] long id)
    {
        var result = await service.DeleteLaboratory(id);
        
        if (!result)
            return NotFound("Laboratory not found");

        return Ok(result);
    }

    [HttpGet("Exists")]
    public async Task<IActionResult> Exists([FromQuery] long id)
    {
        var exists = await service.LaboratoryExists(id);
        return Ok(exists);
    }

    // Invitation endpoints
    [HttpPost("Invite")]
    public async Task<IActionResult> CreateInvitation([FromBody] CreateInvitationDto createInvitationDto)
    {
        // TODO: Get current user ID from authentication context
        long currentUserId = 1; // Placeholder - replace with actual authenticated user ID
        
        var token = await invitationService.CreateInvitation(createInvitationDto, currentUserId);
        var baseUrl = $"{Request.Scheme}://{Request.Host}";
        var invitationUrl = await invitationService.GenerateInvitationUrl(Guid.Parse(token), baseUrl);
        
        return Ok(new { InvitationToken = token, InvitationUrl = invitationUrl });
    }

    [HttpGet("Invitations")]
    public async Task<IActionResult> GetInvitations([FromQuery] long laboratoryId)
    {
        var invitations = await invitationService.GetLaboratoryInvitations(laboratoryId);
        return Ok(invitations);
    }

    [HttpGet("Invitation/{token}")]
    public async Task<IActionResult> GetInvitationDetails(Guid token)
    {
        var invitation = await invitationService.GetInvitationDetails(token);
        
        if (invitation == null)
            return NotFound("Invitation not found");

        return Ok(invitation);
    }

    [HttpPost("AcceptInvitation")]
    public async Task<IActionResult> AcceptInvitation([FromBody] AcceptInvitationDto acceptInvitationDto)
    {
        var result = await invitationService.AcceptInvitation(acceptInvitationDto);
        
        if (!result)
            return BadRequest("Failed to accept invitation or invitation is no longer valid");

        return Ok(result);
    }

    [HttpDelete("CancelInvitation")]
    public async Task<IActionResult> CancelInvitation([FromQuery] long invitationId)
    {
        // TODO: Get current user ID from authentication context
        long currentUserId = 1; // Placeholder - replace with actual authenticated user ID
        
        var result = await invitationService.CancelInvitation(invitationId, currentUserId);
        
        if (!result)
            return NotFound("Invitation not found or you don't have permission to cancel it");

        return Ok(result);
    }

    [HttpGet("PendingInvitations")]
    public async Task<IActionResult> GetPendingInvitations([FromQuery] string email)
    {
        var invitations = await invitationService.GetUserPendingInvitations(email);
        return Ok(invitations);
    }
}
