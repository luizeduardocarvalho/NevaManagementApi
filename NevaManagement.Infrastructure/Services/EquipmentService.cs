
namespace NevaManagement.Infrastructure.Services;

public class EquipmentService : IEquipmentService
{
    private readonly IEquipmentRepository repository;
    private readonly ILocationRepository locationRepository;
    private readonly IMapper mapper;

    public EquipmentService(
        IEquipmentRepository repository,
        ILocationRepository locationRepository,
        IMapper mapper)
    {
        this.repository = repository;
        this.locationRepository = locationRepository;
        this.mapper = mapper;
    }

    public async Task<IList<GetSimpleEquipmentDto>> GetEquipments()
    {
        try
        {
            return await this.repository.GetEquipments();
        }
        catch
        {
            throw new Exception("An error occurred while getting containers.");
        }
    }

    public async Task<bool> AddEquipment(AddEquipmentDto equipmentDto)
    {
        try
        {
            var equipment = mapper.Map<Equipment>(equipmentDto);

            await this.repository.Insert(equipment);
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while creating the equipment.");
        }
    }

    public async Task<bool> EditEquipment(EditEquipmentDto equipmentDto)
    {
        try
        {
            var equipment = await this.repository.GetById(equipmentDto.Id.Value);

            if (equipment is not null)
            {
                if(equipmentDto.LocationId is not null &&
                   await this.locationRepository.GetById(equipmentDto.LocationId.Value) is Location location)
                {
                    equipment.Location = location;
                }

                equipment.Name = equipmentDto.Name;
                equipment.Description = equipmentDto.Description;
                equipment.PropertyNumber = equipmentDto.PropertyNumber;

                return await this.repository.SaveChanges();
            }

            throw new InvalidOperationException($"The equipment with id {equipmentDto.Id} does not exist.");
        }
        catch
        {
            throw new Exception("An error occurred while saving the equipment.");
        }
    }

    public async Task<GetDetailedEquipmentDto> GetDetailedEquipment(long id)
    {
        try
        {
            return await this.repository.GetDetailedEquipment(id);
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment.");
        }
    }
}
