using System.ComponentModel.DataAnnotations;

namespace NevaManagement.Domain.Dtos.Product
{
    public class UseProductDto
    {
        [Required]
        public long ResearcherId { get; set; }

        [Required]
        public long ProductId { get; set; }

        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Quantity should be greater than 0")]
        public double Quantity { get; set; }

        public string Description { get; set; }
    }
}
