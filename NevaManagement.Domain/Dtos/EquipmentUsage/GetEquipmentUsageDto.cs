namespace NevaManagement.Domain.Dtos.EquipmentUsage;

public class GetEquipmentUsageDto
{
    public long Id { get; set; }

    public string Description { get; set; }

    public GetSimpleResearcherDto Researcher { get; set; }

    public DateTimeOffset StartDate { get; set; }

    public DateTimeOffset EndDate { get; set; }
}
