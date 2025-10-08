using NevaManagement.Domain.Dtos.Invitation;

namespace NevaManagement.Infrastructure.Repositories;

public class LaboratoryInvitationRepository : BaseRepository<LaboratoryInvitation>, ILaboratoryInvitationRepository
{
    private readonly NevaManagementDbContext context;

    public LaboratoryInvitationRepository(NevaManagementDbContext context) : base(context)
    {
        this.context = context;
    }

    public async Task<IList<GetInvitationDto>> GetInvitationsByLaboratory(long laboratoryId)
    {
        return await table
            .Where(i => i.LaboratoryId == laboratoryId && i.IsActive)
            .Include(i => i.Laboratory)
            .Include(i => i.InvitedBy)
            .Select(i => new GetInvitationDto
            {
                Id = i.Id,
                LaboratoryId = i.LaboratoryId,
                LaboratoryName = i.Laboratory.Name,
                InviteeEmail = i.InviteeEmail,
                InviteeName = i.InviteeName,
                Role = i.Role,
                CreatedAt = i.CreatedAt,
                ExpiresAt = i.ExpiresAt,
                IsAccepted = i.IsAccepted,
                AcceptedAt = i.AcceptedAt,
                InvitedByName = i.InvitedBy.Name,
                IsExpired = i.ExpiresAt < DateTimeOffset.UtcNow
            })
            .OrderByDescending(i => i.CreatedAt)
            .ToListAsync();
    }

    public async Task<InvitationDetailsDto> GetInvitationByToken(Guid token)
    {
        return await table
            .Where(i => i.InvitationToken == token && i.IsActive)
            .Include(i => i.Laboratory)
            .Include(i => i.InvitedBy)
            .Select(i => new InvitationDetailsDto
            {
                LaboratoryId = i.LaboratoryId,
                LaboratoryName = i.Laboratory.Name,
                LaboratoryDescription = i.Laboratory.Description,
                InviteeName = i.InviteeName,
                Role = i.Role,
                InvitedByName = i.InvitedBy.Name,
                ExpiresAt = i.ExpiresAt,
                IsExpired = i.ExpiresAt < DateTimeOffset.UtcNow,
                IsAccepted = i.IsAccepted
            })
            .FirstOrDefaultAsync();
    }

    public async Task<LaboratoryInvitation> GetInvitationEntityByToken(Guid token)
    {
        return await table
            .Where(i => i.InvitationToken == token && i.IsActive)
            .FirstOrDefaultAsync();
    }

    public async Task<bool> IsInvitationValid(Guid token)
    {
        return await table
            .AnyAsync(i => i.InvitationToken == token 
                        && i.IsActive 
                        && !i.IsAccepted 
                        && i.ExpiresAt > DateTimeOffset.UtcNow);
    }

    public async Task<IList<GetInvitationDto>> GetPendingInvitations(string email)
    {
        return await table
            .Where(i => i.InviteeEmail.ToLower() == email.ToLower() 
                     && i.IsActive 
                     && !i.IsAccepted 
                     && i.ExpiresAt > DateTimeOffset.UtcNow)
            .Include(i => i.Laboratory)
            .Include(i => i.InvitedBy)
            .Select(i => new GetInvitationDto
            {
                Id = i.Id,
                LaboratoryId = i.LaboratoryId,
                LaboratoryName = i.Laboratory.Name,
                InviteeEmail = i.InviteeEmail,
                InviteeName = i.InviteeName,
                Role = i.Role,
                CreatedAt = i.CreatedAt,
                ExpiresAt = i.ExpiresAt,
                IsAccepted = i.IsAccepted,
                AcceptedAt = i.AcceptedAt,
                InvitedByName = i.InvitedBy.Name,
                IsExpired = false
            })
            .OrderByDescending(i => i.CreatedAt)
            .ToListAsync();
    }
}
