namespace NevaManagement.Domain.Dtos.Location
{
    public class GetDetailedLocationDto
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public long? SubLocationId { get; set; }
        public string SubLocationName { get; set; }
    }
}
