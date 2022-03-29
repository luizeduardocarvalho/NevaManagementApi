namespace NevaManagement.Domain.Dtos.Location
{
    public class AddLocationDto
    {
        public string Name { get; set; }
        public string Description { get; set; }
        public long? SublocationId { get; set; }
    }
}
