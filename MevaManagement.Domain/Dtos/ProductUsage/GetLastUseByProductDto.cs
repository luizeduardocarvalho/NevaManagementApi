namespace NevaManagement.Domain.Dtos.ProductUsage;

public class GetLastUseByProductDto
{
    public string ResearcherName { get; set; }
    public DateTimeOffset UseDate { get; set; }
    public double Quantity { get; set; }
    public string Unit { get; set; }
}
