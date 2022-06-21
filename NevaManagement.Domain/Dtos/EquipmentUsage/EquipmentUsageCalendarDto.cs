namespace NevaManagement.Domain.Dtos.EquipmentUsage;

public class EquipmentUsageCalendarDto
{
    public int Month { get; set; }
    public IList<Day> Days { get; set; }
}

public class Day
{
    public int DayNumber { get; set; }
    public IList<Appointment> Appointments { get; set; }
}

public class Appointment
{
    public DateTimeOffset StartDate { get; set; }
    public DateTimeOffset EndDate { get; set; }
}
