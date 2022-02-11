using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace NevaManagement.Domain.Models
{
    [Table("Product")]
    public class Product : BaseEntity
    {
        [MaxLength(100)]
        public string Name { get; set; }

        public string Description { get; set; }

        [ForeignKey("Location_Id")]
        public Location Location { get; set; }

        public double Quantity { get; set; }

        public string Unit { get; set; }
    }
}