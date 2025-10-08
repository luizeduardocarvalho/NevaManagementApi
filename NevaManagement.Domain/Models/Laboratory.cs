namespace NevaManagement.Domain.Models;

[Table("Laboratory")]
public class Laboratory : BaseEntity
{
    [MaxLength(100)]
    public string Name { get; set; }

    public string Description { get; set; }

    [MaxLength(200)]
    public string Address { get; set; }

    [MaxLength(50)]
    public string Phone { get; set; }

    [MaxLength(100)]
    public string Email { get; set; }

    public DateTimeOffset CreatedAt { get; set; }

    public bool IsActive { get; set; }
}