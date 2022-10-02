namespace NevaManagement.Infrastructure.Builders;

public class OrganismBuilder : IBuilder<AddOrganismDto, Organism>
{
    Organism IBuilder<AddOrganismDto, Organism>.Build(AddOrganismDto input)
    {
        return new Organism
        {
            Description = input.Description,
            Name = input.Name,
            CollectionDate = input.CollectionDate.UtcDateTime,
            CollectionLocation = input.CollectionLocation,
            IsolationDate = input.IsolationDate.UtcDateTime,
            OriginPart = input.OriginPart,
            Type = input.Type
        };
    }
}
