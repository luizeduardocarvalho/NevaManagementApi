using System.ComponentModel.DataAnnotations;

namespace NevaManagement.Domain.Dtos.Location
{
    public class EditLocationDto
    {
        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Id should not be 0")]
        public long Id { get; set; }

        [Required]
        public string Name { get; set; }

        public string Description { get; set; }

        public long? SublocationId { get; set; }
    }
}
