namespace NevaManagement.Domain.Dtos.Product;

public class CreateProductDto
{
    public string Description { get; set; }

    [Required]
    public long? LocationId { get; set; }

    [Required]
    public string Name { get; set; }

    public double Quantity { get; set; }

    [Required]
    public string Unit { get; set; }

    public DateTimeOffset ExpirationDate { get; set; }
}
