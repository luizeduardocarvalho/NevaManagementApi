using NevaManagement.Domain.Dtos.Location;

namespace NevaManagement.Domain.Dtos.Product
{
    public class GetDetailedProductDto
    {
        public long Id { get; set; }
        public string Name { get; set; }
        public double Quantity { get; set; }
        public string Unit { get; set; }
        public GetLocationDto Location { get; set; }
    }
}
