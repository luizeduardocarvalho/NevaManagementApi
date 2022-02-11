namespace NevaManagement.Api.Models
{
    public class Location
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public Location SubLocation { get; set; }
    }
}