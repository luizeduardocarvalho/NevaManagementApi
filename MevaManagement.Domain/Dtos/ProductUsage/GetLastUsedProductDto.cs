namespace NevaManagement.Domain.Dtos.ProductUsage
{
    public class GetLastUsedProductDto
    {
        public long ProductId { get; set; }
        public string ProductName { get; set; }
        public string LocationName { get; set; }
        public double Quantity { get; set; }
        public string Unit { get; set; }
    }
}
