using System.ComponentModel.DataAnnotations;

namespace NevaManagement.Domain.Dtos.Product
{
    public class EditProductDto
    {
        [Required]
        public long? Id { get; set; }

        public string Name { get; set; }

        public long? LocationId { get; set; }

        public string Description { get; set; }
    }
}
