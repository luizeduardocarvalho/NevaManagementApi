using System.ComponentModel.DataAnnotations.Schema;

namespace NevaManagement.Domain.Models
{
    [Table("Location")]
    public class Location : BaseEntity
    {
        public string Name { get; set; }

        public string Description { get; set; }

        [ForeignKey("Sub_Location_Id")]
        public Location SubLocation { get; set; }
    }
}