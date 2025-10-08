using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Infrastructure.Repositories;

public class LaboratoryRepository : BaseRepository<Laboratory>, ILaboratoryRepository
{
    public LaboratoryRepository(NevaManagementDbContext context) : base(context)
    {
    }

    public async Task<IList<GetLaboratoryDto>> GetAllLaboratories()
    {
        return await table
            .Where(l => l.IsActive)
            .Select(l => new GetLaboratoryDto
            {
                Id = l.Id,
                Name = l.Name,
                Description = l.Description,
                Address = l.Address,
                Phone = l.Phone,
                Email = l.Email,
                CreatedAt = l.CreatedAt,
                IsActive = l.IsActive
            })
            .ToListAsync();
    }

    public async Task<IList<GetSimpleLaboratoryDto>> GetSimpleLaboratories()
    {
        return await table
            .Where(l => l.IsActive)
            .Select(l => new GetSimpleLaboratoryDto
            {
                Id = l.Id,
                Name = l.Name,
                IsActive = l.IsActive
            })
            .ToListAsync();
    }

    public async Task<GetLaboratoryDto> GetLaboratoryById(long id)
    {
        return await table
            .Where(l => l.Id == id)
            .Select(l => new GetLaboratoryDto
            {
                Id = l.Id,
                Name = l.Name,
                Description = l.Description,
                Address = l.Address,
                Phone = l.Phone,
                Email = l.Email,
                CreatedAt = l.CreatedAt,
                IsActive = l.IsActive
            })
            .FirstOrDefaultAsync();
    }

    public async Task<bool> LaboratoryExists(long id)
    {
        return await table.AnyAsync(l => l.Id == id && l.IsActive);
    }
}
