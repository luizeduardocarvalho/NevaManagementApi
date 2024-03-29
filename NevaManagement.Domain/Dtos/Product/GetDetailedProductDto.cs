﻿namespace NevaManagement.Domain.Dtos.Product;

public class GetDetailedProductDto
{
    public long Id { get; set; }
    
    public string Name { get; set; }
    
    public double Quantity { get; set; }
    
    public double QuantityUsedInTheLastThreeMonths { get; set; }
    
    public string Description { get; set; }
    
    public string Unit { get; set; }
    
    public string Formula { get; set; }
    
    public GetLocationDto Location { get; set; }
    
    public DateTimeOffset ExpirationDate { get; set; }
}
