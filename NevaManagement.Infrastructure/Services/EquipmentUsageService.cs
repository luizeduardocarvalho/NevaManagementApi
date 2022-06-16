namespace NevaManagement.Infrastructure.Services;

public class EquipmentUsageService : IEquipmentUsageService
{
    private readonly IEquipmentUsageRepository repository;

    public EquipmentUsageService(IEquipmentUsageRepository repository)
    {
        this.repository = repository;
    }

    public async Task<dynamic> GetEquipmentUsageCalendar(long equipmentId)
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
                select new
                {
                    Month = all.Key,
                    Days = all.Select(x =>
                      new
                      {
                          Day = x.Key,
                          Times = x.Select(t => new { t.StartDate, t.EndDate })
                      })
                };

            return nestedDates;
        }
        catch
        {
            throw new Exception("An error occurred while getting the equipment usage calendar");
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
