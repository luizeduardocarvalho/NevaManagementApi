using System;

namespace NevaManagement.Domain.Dtos.Container
{
    public class GetDetailedContainerDto
    {
        public string Name { get; set; }
        public string CultureMedia { get; set; }
        public string ResearcherName { get; set; }
        public string Description { get; set; }
        public string OriginName { get; set; }
        public DateTimeOffset CreationDate { get; set; }
    }
}
