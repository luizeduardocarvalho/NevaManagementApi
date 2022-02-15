namespace NevaManagement.Domain.Dtos.Product
{
    public class CreateProductDto
    {
        public string Description { get; set; }
        public long LocationId { get; set; }
        public string Name { get; set; }
        public double Quantity { get; set; }
        public string Unit { get; set; }
    }
}
