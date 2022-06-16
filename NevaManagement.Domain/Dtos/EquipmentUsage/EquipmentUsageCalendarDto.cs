namespace NevaManagement.Domain.Dtos.EquipmentUsage;

public class EquipmentUsageCalendarDto
{
    public int Month { get; set; }
    public Days Days { get; set; }
}

public class Days
{
    public int Day { get; set; }
    public Times Times { get; set; }
}

public class Times
{
    public DateTimeOffset StartDate { get; set; }
    public DateTimeOffset EndDate { get; set; }
}

//var expectedResult = new
//{
//    Month = 1,
//    Days = new
//    {
//        Day = 1,
//        Times = new
//        {
//            StartDate = new DateTime(2022, 1, 1, 10, 30, 0),
//            EndDate = new DateTime(2022, 1, 1, 11, 30, 0)
//        }
//    }
//};
