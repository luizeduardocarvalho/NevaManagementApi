namespace NevaManagement.Domain.Dtos.Equipment;

public class GetDetailedEquipmentDto
{
    public long Id { get; set; }

    public string Name { get; set; }

    public string Description { get; set; }

    public string PropertyNumber { get; set; }

    public GetLocationDto Location { get; set; }

    public IList<GetEquipmentUsageDto> UsageList { get; set; }
}
