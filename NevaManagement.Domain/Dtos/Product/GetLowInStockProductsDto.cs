namespace NevaManagement.Domain.Dtos.Product;

public class GetLowInStockProductsDto
{
    public long Id { get; set; }

    public string Name { get; set; }

    public string Unit { get; set; }

    public double Quantity { get; set; }
}
