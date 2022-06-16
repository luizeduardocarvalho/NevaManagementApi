using NevaManagement.Domain.Dtos.EquipmentUsage;

namespace NevaManagement.Domain.Dtos.Equipment;

public class GetDetailedEquipmentDto
{
    public long Id { get; set; }

    public string Name { get; set; }

    public string Description { get; set; }

    public IList<GetEquipmentUsageDto> UsageList { get; set; }
}
