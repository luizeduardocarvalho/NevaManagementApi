namespace NevaManagement.Api.Models
{
    public class Product
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public Location Location { get; set; }
        public double Quantity { get; set; }
        public string Unit { get; set; }
    }
}