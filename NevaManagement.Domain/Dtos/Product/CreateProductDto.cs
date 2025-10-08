using System.Globalization;
using System.Text.Json.Serialization;

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

    public string ExpirationDate { get; set; }

    [Required]
    public long LaboratoryId { get; set; }

    [JsonIgnore]
    public DateTimeOffset? ParsedExpirationDate
    {
        get 
        {
            if(!string.IsNullOrEmpty(ExpirationDate))
            {
                return DateTimeOffset.Parse(ExpirationDate, new CultureInfo("en-US")); 
            }
            else
            {
                return null;
            }
        }
    }
}
