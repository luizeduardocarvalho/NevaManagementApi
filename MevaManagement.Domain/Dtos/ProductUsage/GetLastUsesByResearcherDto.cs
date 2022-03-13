using System;

namespace NevaManagement.Domain.Dtos.ProductUsage
{
    public class GetLastUsesByResearcherDto
    {
        public string ResearcherName { get; set; }
        public DateTimeOffset UsageDate { get; set; }
        public double Quantity { get; set; }
        public string Unit { get; set; }
    }
}
