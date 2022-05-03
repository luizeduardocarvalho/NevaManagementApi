using System;
using System.ComponentModel.DataAnnotations;

namespace NevaManagement.Domain.Dtos.Organism
{
    public class AddOrganismDto
    {
        [Required]
        public string Name { get; set; }

        public string Type { get; set; }

        public string Description { get; set; }

        public DateTimeOffset CollectionDate { get; set; }

        public string CollectionLocation { get; set; }

        public DateTimeOffset IsolationDate { get; set; }

        public long? OriginOrganismId { get; set; }

        public string OriginPart { get; set; }
    }
}
