namespace NevaManagement.Infrastructure.Services;

public class EquipmentUsageService : IEquipmentUsageService
{
    private readonly IEquipmentUsageRepository repository;

    public EquipmentUsageService(IEquipmentUsageRepository repository)
    {
        this.repository = repository;
    }

    public async Task<IList<EquipmentUsageCalendarDto>> GetEquipmentUsageCalendar(long equipmentId)
    {
        try
        {
            var result = await this.repository.GetEquipmentUsageByEquipment(equipmentId);

            var nestedDates =
                from month in result
                group month by month.StartDate.Month into months
                from days in (
                    from day in months
                    group day by day.StartDate.Day
                )
                group days by months.Key into all
                select new EquipmentUsageCalendarDto
                {
                    Month = all.Key,
                    Days = all.Select(x =>
                      new Day
                      {
                          DayNumber = x.Key,
                          Appointments = x.Select(t => 
                            new Appointment 
                            { 
                                StartDate = t.StartDate,
                                EndDate = t.EndDate 
                            }).ToList()
                      }).ToList()
                };

            return nestedDates.ToList();
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment usage calendar.");
        }
    }

    public async Task<IList<GetEquipmentUsageDto>> GetEquipmentUsageHistory(long id)
    {
        try
        {
            return await this.repository.GetEquipmentUsageByEquipment(id);
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
                    useEquipmentDto.StartDate);

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
                StartDate = useEquipmentDto.StartDate,
                EndDate = useEquipmentDto.EndDate
            };

            await this.repository.Insert(equipmentUsage);
            return await this.repository.SaveChanges();
        }
        catch
        {
            throw new Exception("An error occurred while trying to use the equipment");
        }
    }
}
