using NevaManagement.Domain.Dtos.Laboratory;

namespace NevaManagement.Infrastructure.Services;

public class LaboratoryService : ILaboratoryService
{
    private readonly ILaboratoryRepository repository;

    public LaboratoryService(ILaboratoryRepository repository)
    {
        this.repository = repository;
    }

    public async Task<IList<GetLaboratoryDto>> GetAllLaboratories()
    {
        try
        {
            return await repository.GetAllLaboratories();
        }
        catch
        {
            throw new Exception("An error occurred while getting laboratories.");
        }
    }

    public async Task<IList<GetSimpleLaboratoryDto>> GetSimpleLaboratories()
    {
        try
        {
            return await repository.GetSimpleLaboratories();
        }
        catch
        {
            throw new Exception("An error occurred while getting laboratories.");
        }
    }

    public async Task<GetLaboratoryDto> GetLaboratoryById(long id)
    {
        try
        {
            return await repository.GetLaboratoryById(id);
        }
        catch
        {
            throw new Exception("An error occurred while getting the laboratory.");
        }
    }

    public async Task<bool> CreateLaboratory(CreateLaboratoryDto createLaboratoryDto)
    {
        try
        {
            var laboratory = new Laboratory
            {
                Name = createLaboratoryDto.Name,
                Description = createLaboratoryDto.Description,
                Address = createLaboratoryDto.Address,
                Phone = createLaboratoryDto.Phone,
                Email = createLaboratoryDto.Email,
                CreatedAt = DateTimeOffset.UtcNow,
                IsActive = true
            };

            await repository.Insert(laboratory);
            return await repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while creating the laboratory.");
        }
    }

    public async Task<bool> UpdateLaboratory(EditLaboratoryDto editLaboratoryDto)
    {
        try
        {
            var laboratory = await repository.GetById(editLaboratoryDto.Id);
            if (laboratory == null)
                return false;

            laboratory.Name = editLaboratoryDto.Name;
            laboratory.Description = editLaboratoryDto.Description;
            laboratory.Address = editLaboratoryDto.Address;
            laboratory.Phone = editLaboratoryDto.Phone;
            laboratory.Email = editLaboratoryDto.Email;
            laboratory.IsActive = editLaboratoryDto.IsActive;

            return await repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while updating the laboratory.");
        }
    }

    public async Task<bool> DeleteLaboratory(long id)
    {
        try
        {
            var laboratory = await repository.GetById(id);
            if (laboratory == null)
                return false;

            laboratory.IsActive = false;
            return await repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while deleting the laboratory.");
        }
    }

    public async Task<bool> LaboratoryExists(long id)
    {
        try
        {
            return await repository.LaboratoryExists(id);
        }
        catch
        {
            throw new Exception("An error occurred while checking laboratory existence.");
        }
    }
}
