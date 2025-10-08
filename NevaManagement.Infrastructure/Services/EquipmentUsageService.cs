namespace NevaManagement.Infrastructure.Services;

public class EquipmentUsageService : IEquipmentUsageService
{
    private readonly IEquipmentUsageRepository repository;

    public EquipmentUsageService(IEquipmentUsageRepository repository)
    {
        this.repository = repository;
    }

    public async Task<IList<EquipmentUsageCalendarDto>> GetEquipmentUsageCalendar(long equipmentId, long laboratoryId)
    {
        try
        {
            var result = await this.repository.GetEquipmentUsageByEquipment(equipmentId, laboratoryId);

            var nestedDates = result
                .GroupBy(x => x.StartDate.Month)
                .Select(monthGroup => new EquipmentUsageCalendarDto
                {
                    Month = monthGroup.Key,
                    Days = monthGroup
                        .GroupBy(x => x.StartDate.Day)
                        .Select(dayGroup => new Day
                        {
                            DayNumber = dayGroup.Key,
                            Appointments = dayGroup
                                .Select(usage => new Appointment
                                {
                                    StartDate = usage.StartDate,
                                    EndDate = usage.EndDate
                                })
                                .ToList()
                        })
                        .ToList()
                });

            return nestedDates.ToList();
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment usage calendar.");
        }
    }

    public async Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageHistory(long id, long laboratoryId)
    {
        try
        {
            return await this.repository.GetEquipmentUsageByEquipment(id, laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment usage history.");
        }
    }

    public async Task<bool> UseEquipment(UseEquipmentDto useEquipmentDto)
    {
        try
        {
            var dates = await this.repository
                .GetEquipmentUsageByDay(
                    useEquipmentDto.EquipmentId,
                    useEquipmentDto.StartDate,
                    useEquipmentDto.LaboratoryId);

            foreach (var date in dates)
            {
                if(useEquipmentDto.EndDate > date.StartDate && useEquipmentDto.EndDate < date.EndDate ||
                    useEquipmentDto.StartDate > date.StartDate && useEquipmentDto.StartDate < date.EndDate ||
                    useEquipmentDto.StartDate < date.StartDate && useEquipmentDto.EndDate > date.EndDate)
                {
                    return false;
                }
            }

            var equipmentUsage = new EquipmentUsage()
            {
                ResearcherId = useEquipmentDto.ResearcherId,
                EquipmentId = useEquipmentDto.EquipmentId,
                Description = useEquipmentDto.Description,
                StartDate = useEquipmentDto.StartDate.UtcDateTime,
                EndDate = useEquipmentDto.EndDate.UtcDateTime,
                LaboratoryId = useEquipmentDto.LaboratoryId
            };

            await this.repository.Insert(equipmentUsage);
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while trying to use the equipment");
        }
    }

    async Task<IList<GetEquipmentUsageDto>> IEquipmentUsageService.GetEquipmentUsage(long equipmentid, int page, long laboratoryId)
    {
        try
        {
            return await this.repository.GetEquipmentUsage(equipmentid, page, laboratoryId);
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment history.");
        }
    }
}
